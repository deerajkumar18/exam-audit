# exam-audit

About :

The application is designed to identify and report suspicious answering pattern amongst the exam candidates . To generate the the suspicion report , the edit pattern along with the submission time is used to generate correlation score between two candidates . This approach differs from the other online proctored exams as they generally flag candidates real time based on the suspicious activities detected from the webcam or by monitoring the screen . 

Why Hyperledger-fabric ?
This is a permissioned private blockchain which can be accessed by client only using the fabric identity . Also , it allows us to query the trnasaction history for a key . By leveraging this , using questionid~studentid~examid as key , answers edited by the candidates can be queried and analyse the results . This ledger technology can be used in place of traditional DB to retrive edits made on a question and compare across candidates whereas the traditional DB stores only the world state .

Pre-requisite:
1. Install fabric-samples (used to create local fabric network and keys to autenticate)
2. Install Go
3. Install docker and docker-compose 

Steps to execute locally :
1. Bring up the hyperledger fabric network
2. Execute the client application
3. Submit the answer using the /submit-answer api :
    curl -X POST http://localhost:8080/submit-answer  -H "Content-Type: application/json"  -d '{"examId": "exam123","questionId": "Q1","ans": "Option B","StudentID":"s1"}'
4. Request for audit report using the /audit-report api :
     curl -v 'http://localhost:8080/audit-answer?examID=exam123&instructorId=i1'

Sequence Diagram :
1. Answers submission:
<img width="501" height="372" alt="submit answer sequence" src="https://github.com/user-attachments/assets/3fdbdd7f-53cf-4653-b5d4-98bccd544560" />

2. Audit report generation:
<img width="381" height="452" alt="generate audit report" src="https://github.com/user-attachments/assets/914ac54c-1488-4633-9438-00ca2c2922dc" />
