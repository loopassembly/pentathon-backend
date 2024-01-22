if (isTeam) {
        // Construct data array based on form input for the team
        let teamData = [
          getValue('#teamName'),
          getValue('#teamLeadName'),
          getValue('#leadRegisterNumber'),
          getValue('#leadSRMEmail'),
          getValue('#leadYearOfStudy'),
          getValue('#leadDepartment'),
          getValue('#leadSection'),
          getValue('#leadWhatsappNo'),
          getValue('#leadFAName'),
          getValue('#ally1Name'),
          getValue('#ally1RegisterNumber'),
          getValue('#ally1SRMEmail'),
          getValue('#ally1YearOfStudy'),
          getValue('#ally1Department'),
          getValue('#ally1Section'),
          getValue('#ally1WhatsappNo'),
          getValue('#ally1FAName'),
         
          // ...
          getValue('#transactionNumber'),
          getValue('#paymentScreenshot') 
        ];

        // Construct the request body
        let requestBody = {
          "data": [teamData]
        };

        // Make the fetch request
        fetch("http://127.0.0.1:3000/teamhandler", {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json'
          },
          body: JSON.stringify(requestBody),
          redirect: 'follow'
        })
          .then(response => response.text())
          .then(result => console.log(result))
          .catch(error => console.log('error', error));
      } else {
        // Handle solo form submission similarly if needed
      }
    });