# maths

A web application for kids to practice mental math and column form operations.

Initially developed in Python / Flask / SQLite - meant to be deployed on my old Synology NAS - in 2017, only for my kids, I'm now porting it to Golang / MongoDB and adding new features to make it easier to install and administrate.

Amongst new features, there are now 3 roles:
- Admins have access to the user management UI,
- Teachers / Parents have access to the homework management UI where they can configure and assign exercises,
- Students can manage their own profile and do the homework assigned by their teachers/parents.

Unchanges features include:
- Choice between mental math and column form operations.
    - Mental math: do the predefined number of operations in 10 minutes, no errors allowed.
    - Column form operations: do the predefined number of operations in 30 minutes, errors are allowed.
- UI optimized for tablets/touch screens, based on JQuery, Bootstrap and WebSocket.
- Summary displayed at the end of the exercise
- History of exercises done recently can be displayed, with some filters
