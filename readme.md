https://github.com/blessingman/education-platform/blob/main/internal/bot/bot.go



TODO:


	
// TODO: stplit this logic and move to code genarator serveice 
// && and 
// 1. create the method (ver code rep ) for return the diff between create_at and expiresAt
// 2. Inject this repo to code gen service 
// 3. In service add method IsValid for verify otp valid or not 
// 
// 1. review  emailService 
// 2. create email client && email client interface with one method: send 
// 3. verify of creds and of conntect to smtp for example: test request 
// 4. create the email formatter struct and move to those place the logic FormatBSUEmail




Mail for check otp password:

https://webmail.bsu.by/owa/#path=/mail



task: 

 create validotor service and other email prefixes



18.04.2025

task: 
 - [ ] implement the logic of fsm provider 


 fms provider
 tx begix -> 
 - start of commit to fsm state table 
 - start of commit verification codes table 


This class is implement the logic of two repository:
 - verification code repo
 - fsm repo 



 - [ ] telegram hahdler decorator and decorator as variafic function for implement the middlware logic


