const sheetName = 'pentathon';
const scriptProp = PropertiesService.getScriptProperties();

function intialSetup() {
  const activeSpreadsheet = SpreadsheetApp.getActiveSpreadsheet();
  scriptProp.setProperty('key', activeSpreadsheet.getId());
}

// function emailExists(sheet, email) {
//   const emailColumnIndex = 3; // Assuming email is in the 3rd column, adjust as needed
//   const emails = sheet.getRange(2, emailColumnIndex, sheet.getLastRow() - 1, 1).getValues().flat();
//   return emails.includes(email);
// }

function emailExists(sheet, email) {
  const emailColumnIndex = 4; // Assuming email is in the 4th column, adjust as needed
  const emails = sheet.getRange(2, emailColumnIndex, sheet.getLastRow() - 1, 1).getValues().flat();
  return emails.some(existingEmail => existingEmail.trim().toLowerCase() === email.trim().toLowerCase());
}

function doPost(e) {
  const lock = LockService.getScriptLock();
  lock.tryLock(10000);

  try {
    const doc = SpreadsheetApp.openById(scriptProp.getProperty('key'));
    const sheet = doc.getSheetByName(sheetName);

    const headers = sheet.getRange(1, 1, 1, sheet.getLastColumn()).getValues()[0];
    const email = e.parameter['SRMIST e-mail'];

    // Check if email already exists
    if (emailExists(sheet, email)) {
      return ContentService
        .createTextOutput(JSON.stringify({ 'result': 'error', 'error': 'Email already exists' }))
        .setMimeType(ContentService.MimeType.JSON);
    }

    const nextRow = sheet.getLastRow() + 1;

    const newRow = headers.map(function(header) {
      return header === 'Date' ? new Date() : e.parameter[header];
    });

    sheet.getRange(nextRow, 1, 1, newRow.length).setValues([newRow]);

    return ContentService
      .createTextOutput(JSON.stringify({ 'result': 'success', 'row': nextRow }))
      .setMimeType(ContentService.MimeType.JSON);
  } catch (e) {
    console.error("Error:", e);
    return ContentService
      .createTextOutput(JSON.stringify({ 'result': 'error', 'error': e }))
      .setMimeType(ContentService.MimeType.JSON);
  } finally {
    lock.releaseLock();
  }
}
