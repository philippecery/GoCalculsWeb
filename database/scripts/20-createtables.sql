USE maths;

/********** CONSTANTS **********/

/* UserRole table */
DROP TABLE IF EXISTS UserRole;
CREATE TABLE UserRole (
    RoleID      TINYINT(1) UNSIGNED NOT NULL PRIMARY KEY,
    RoleI18n    VARCHAR(50) NOT NULL
);

INSERT INTO UserRole (RoleID, RoleI18n) values
    (0, 'Child/Student'),
    (1, 'Parent/Teacher'),
    (2, 'Principal'),
    (3, 'Admin');


/* UserStatus table */
DROP TABLE IF EXISTS UserStatus;
CREATE TABLE UserStatus (
    StatusID      TINYINT(1) UNSIGNED NOT NULL PRIMARY KEY,
    StatusI18n    VARCHAR(50) NOT NULL
);

INSERT INTO UserStatus (StatusID, StatusI18n) values
    (0, 'Unregistered'),
    (1, 'Disabled'),
    (2, 'Enabled');


/* TeamType table */
DROP TABLE IF EXISTS TeamType;
CREATE TABLE TeamType (
    TypeID      TINYINT(1) UNSIGNED NOT NULL PRIMARY KEY,
    TypeI18n    VARCHAR(50) NOT NULL
);

INSERT INTO TeamType (TypeID, TypeI18n) values
    (1, 'Family'),
    (2, 'School');


/* TeamStatus table */
DROP TABLE IF EXISTS TeamStatus;
CREATE TABLE TeamStatus (
    StatusID      TINYINT(1) UNSIGNED NOT NULL PRIMARY KEY,
    StatusI18n    VARCHAR(50) NOT NULL
);

INSERT INTO TeamStatus (StatusID, StatusI18n) values
    (0, 'Unregistered'),
    (1, 'Disabled'),
    (2, 'Enabled');

/* TeamUserRole table */
DROP TABLE IF EXISTS TeamUserRole;
CREATE TABLE TeamUserRole (
    RoleID      TINYINT(1) UNSIGNED NOT NULL PRIMARY KEY,
    RoleI18n    VARCHAR(50) NOT NULL
);

INSERT INTO TeamUserRole (RoleID, RoleI18n) values
    (0, 'Normal'),
    (1, 'Admin'),
    (2, 'Root');

/* HomeworkStatus table */
DROP TABLE IF EXISTS HomeworkStatus;
CREATE TABLE HomeworkStatus (
    StatusID      TINYINT(1) UNSIGNED NOT NULL PRIMARY KEY,
    StatusI18n    VARCHAR(50) NOT NULL
);

INSERT INTO HomeworkStatus (StatusID, StatusI18n) values
    (0, 'Draft'),
    (1, 'Online'),
    (2, 'Archived');

/* SessionStatus table */
DROP TABLE IF EXISTS SessionStatus;
CREATE TABLE SessionStatus (
    StatusID      TINYINT(1) UNSIGNED NOT NULL PRIMARY KEY,
    StatusI18n    VARCHAR(50) NOT NULL
);

INSERT INTO SessionStatus (StatusID, StatusI18n) values
    (0, 'Cancel'),
    (1, 'Failed'),
    (2, 'Timeout'),
    (3, 'Success');

/********** ADMINISTRATION **********/

/* OperandRanges table */
DROP TABLE IF EXISTS OperandRanges;
CREATE TABLE OperandRanges (
    OperandRangeID      TINYINT(2) UNSIGNED NOT NULL PRIMARY KEY,
    Operand1RangeMin    INT NOT NULL,
    Operand1RangeMax    INT NOT NULL,
    Operand1DecimalMax  TINYINT(1) UNSIGNED NOT NULL,
    Operand2RangeMin    INT NOT NULL,
    Operand2RangeMax    INT NOT NULL,
    Operand2DecimalMax  TINYINT(1) UNSIGNED NOT NULL
);

INSERT INTO OperandRanges (OperandRangeID, Operand1RangeMin, Operand1RangeMax, Operand1DecimalMax, Operand2RangeMin, Operand2RangeMax, Operand2DecimalMax) values
    (11, 1, 100, 0, 1, 10, 0),
    (12, 10, 100, 0, 1, 10, 0),
    (13, 2, 10, 0, 2, 10, 0),
    (14, 10, 100, 0, 2, 10, 0),
    (21, 100, 1000000, 2, 100, 100000, 2),
    (22, 100000, 1000000, 2, 1000, 100000, 2),
    (23, 100, 100000, 2, 100, 100000, 2),
    (24, 1000, 100000, 0, 10, 1000, 0);


/* GradeHomeworkTypes table */
DROP TABLE IF EXISTS HomeworkTypes;
CREATE TABLE HomeworkTypes (
    TypeID                  TINYINT(1) UNSIGNED NOT NULL PRIMARY KEY,
	TypeI18n                VARCHAR(50) NOT NULL,
	TypeLogo                VARCHAR(50) NOT NULL,
	AdditionRangeID         TINYINT(2) UNSIGNED NOT NULL,
	SubstractionRangeID     TINYINT(2) UNSIGNED NOT NULL,
	MultiplicationRangeID   TINYINT(2) UNSIGNED NOT NULL,
	DivisionRangeID         TINYINT(2) UNSIGNED NOT NULL,
    FOREIGN KEY (AdditionRangeID)       REFERENCES OperandRanges (OperandRangeID),
    FOREIGN KEY (SubstractionRangeID)   REFERENCES OperandRanges (OperandRangeID),
    FOREIGN KEY (MultiplicationRangeID) REFERENCES OperandRanges (OperandRangeID),
    FOREIGN KEY (DivisionRangeID)       REFERENCES OperandRanges (OperandRangeID)
);

INSERT INTO HomeworkTypes (TypeID, TypeI18n, TypeLogo, AdditionRangeID, SubstractionRangeID, MultiplicationRangeID, DivisionRangeID) values
    (1, 'mentalmath', 'hourglass', 11, 12, 13, 14),
    (2, 'columnform', 'pencil', 21, 22, 23, 24);


/* Teams table */
DROP TABLE IF EXISTS Teams;
CREATE TABLE Teams (
    TeamID           CHAR(64) NOT NULL PRIMARY KEY,
    TeamTypeID       TINYINT(1) UNSIGNED NOT NULL,
    TeamName         VARCHAR(50),
    TeamLanguageCode CHAR(5) NOT NULL,
    TeamStatusID     TINYINT(1) UNSIGNED NOT NULL,
    FOREIGN KEY (TeamTypeID)     REFERENCES TeamType (TypeID),
    FOREIGN KEY (TeamStatusID)   REFERENCES TeamStatus (StatusID)
);


/* Users table */
DROP TABLE IF EXISTS Users;
CREATE TABLE Users (
    UserID          CHAR(64) NOT NULL PRIMARY KEY,
    UserPassword    VARCHAR(50),
    PasswordDate    DATETIME,
    FullName        VARCHAR(260),
    LanguageCode    CHAR(5) NOT NULL,
    EmailAddress    VARCHAR(260) NOT NULL,
    EmailAddressTmp VARCHAR(260),
    RoleID          TINYINT(1) UNSIGNED NOT NULL,
    StatusID        TINYINT(1) UNSIGNED NOT NULL,
	LastConnection  DATETIME,
	Token           VARCHAR(50),
	Expires         DATETIME,
	FailedAttempts  TINYINT(1) UNSIGNED,
    FOREIGN KEY (RoleID)    REFERENCES UserRole (RoleID),
    FOREIGN KEY (StatusID)  REFERENCES UserStatus (StatusID)
);


/* TeamsUsers table */
DROP TABLE IF EXISTS TeamsUsers;
CREATE TABLE TeamsUsers (
    TeamID      CHAR(64) NOT NULL,
    UserID      CHAR(64) NOT NULL,
    RoleID      TINYINT(1) UNSIGNED NOT NULL,
    TeamDefault TINYINT(1) UNSIGNED NOT NULL,
    PRIMARY KEY (TeamID, UserID),
    FOREIGN KEY (TeamID)     REFERENCES Teams (TeamID) ON DELETE CASCADE,
    FOREIGN KEY (UserID)    REFERENCES Users (UserID) ON DELETE CASCADE,
    FOREIGN KEY (RoleID)    REFERENCES TeamUserRole (RoleID)
);


/* Grades table */
DROP TABLE IF EXISTS Grades;
CREATE TABLE Grades (
    GradeID             CHAR(64) NOT NULL PRIMARY KEY,
    OwnerID             CHAR(64) NOT NULL,
    GradeName           VARCHAR(260) NOT NULL,
    GradeDescription    VARCHAR(15000),
    FOREIGN KEY (OwnerID) REFERENCES Users (UserID)
);


/* Homeworks table */
DROP TABLE IF EXISTS Homeworks;
CREATE TABLE Homeworks (
    HomeworkID          CHAR(64) NOT NULL PRIMARY KEY,
    HomeworkName        VARCHAR(260) NOT NULL,
    GradeID             CHAR(64) NOT NULL,
    HomeworkTypeID      TINYINT(1) UNSIGNED NOT NULL,
	NbAdditions         TINYINT UNSIGNED,
	NbSubstractions     TINYINT UNSIGNED,
	NbMultiplications   TINYINT UNSIGNED,
	NbDivisions         TINYINT UNSIGNED,
	HomeworkTime        TINYINT UNSIGNED,
    StatusID            TINYINT(1) UNSIGNED NOT NULL,
    -- UNIQUE KEY (HomeworkTypeID, NbAdditions, NbSubstractions, NbMultiplications, NbDivisions, HomeworkTime),
    FOREIGN KEY (HomeworkTypeID)    REFERENCES HomeworkTypes (TypeID),
    FOREIGN KEY (GradeID)           REFERENCES Grades (GradeID),
    FOREIGN KEY (StatusID)          REFERENCES HomeworkStatus (StatusID)
);


/* Users table */
DROP TABLE IF EXISTS Students;
CREATE TABLE Students (
    UserID  CHAR(64) NOT NULL,
	GradeID CHAR(64) NOT NULL,
    PRIMARY KEY (UserID /*, GradeID*/),
    FOREIGN KEY (UserID)    REFERENCES Users (UserID) ON DELETE CASCADE,
    FOREIGN KEY (GradeID)   REFERENCES Grades (GradeID)
);


/********** HOMEWORK DATA **********/

/* HomeworkSessions table */
DROP TABLE IF EXISTS HomeworkSessions;
CREATE TABLE HomeworkSessions (
    SessionID   CHAR(64) NOT NULL PRIMARY KEY,
	UserID      CHAR(64) NOT NULL,
	StartTime   DATETIME NOT NULL,
    EndTime     DATETIME,
	HomeworkID  CHAR(64) NOT NULL,
	StatusID    TINYINT(1) UNSIGNED NOT NULL,
    FOREIGN KEY (UserID)     REFERENCES Users (UserID),
    FOREIGN KEY (HomeworkID) REFERENCES Homeworks (HomeworkID),
    FOREIGN KEY (StatusID)   REFERENCES SessionStatus (StatusID)
);


/* SessionOperations table */
DROP TABLE IF EXISTS SessionOperations;
CREATE TABLE SessionOperations (
    SessionID   CHAR(64) NOT NULL,
    OperationID SMALLINT UNSIGNED NOT NULL,
	OperatorID  TINYINT(1) UNSIGNED NOT NULL,
	Operand1    INT NOT NULL,
	Operand2    INT NOT NULL,
	Answer      INT,
	Answer2     INT,
	StatusID    TINYINT(1) NOT NULL,
    PRIMARY KEY (SessionID, OperationID),
    FOREIGN KEY (SessionID) REFERENCES HomeworkSessions (SessionID)
);
