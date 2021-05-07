package dataaccess

import (
	"context"
	"database/sql"
	"log"

	"github.com/philippecery/maths/webapp/database"
	"github.com/philippecery/maths/webapp/database/model"
)

const queryGetTeamsByUserID = `
SELECT
	o.TeamID, o.TeamTypeID, o.TeamName, o.TeamLanguageCode, o.TeamStatusID
FROM Users u
	INNER JOIN TeamsUsers ou ON u.UserID = ou.UserID
	INNER JOIN Teams o ON ou.TeamID = o.TeamID
WHERE u.UserID = ?
`

// GetTeamsByUserID returns the organizations for the given userID
func GetTeamsByUserID(userID string) []*model.Team {
	var teams []*model.Team
	if rows, err := database.Connection.QueryContext(context.TODO(), queryGetTeamsByUserID, userID); err == nil {
		defer rows.Close()
		for rows.Next() {
			team := new(model.Team)
			if err = rows.Scan(team.TeamID, team.Type, team.Name, team.Language, team.Status); err == nil {
				teams = append(teams, team)
			}
		}
	}
	return teams
}

const queryGetDefaultTeamByUserID = `
SELECT
	o.TeamID, o.TeamTypeID, o.TeamName, o.TeamLanguageCode, o.TeamStatusID
FROM Users u
	INNER JOIN TeamsUsers ou ON u.UserID = ou.UserID
	INNER JOIN Teams o ON ou.TeamID = o.TeamID
WHERE u.UserID = ? AND ou.TeamDefault = ?
`

// GetDefaultTeamByUserID returns the organizations for the given userID
func GetDefaultTeamByUserID(userID string) *model.Team {
	team := new(model.Team)
	err := database.Connection.QueryRowContext(context.TODO(), queryGetDefaultTeamByUserID, userID, true).Scan(team.TeamID, team.Type, team.Name, team.Language, team.Status)
	switch {
	case err == sql.ErrNoRows:
		log.Printf("database: user %s doesn't exist", userID)
	case err != nil:
		log.Printf("database: unable to find user %s. Cause: %v", userID, err)
		return nil
	}
	return team
}

const (
	queryUnsetDefaultTeam = "UPDATE TeamsUsers SET TeamDefault = ? WHERE UserID = ? AND TeamDefault = ?"
	querySetDefaultTeam   = "UPDATE TeamsUsers SET TeamDefault = ? WHERE UserID = ? AND TeamID = ?"
)

// UpdateDefaultTeam updates the number of failed login attempts for the User model where userid field is the provided id
func UpdateDefaultTeam(userID, teamID string) error {
	var err error
	var tx *sql.Tx
	ctx := context.TODO()
	if tx, err = database.Connection.BeginTx(ctx, nil); err == nil {
		if _, err = tx.ExecContext(ctx, queryUnsetDefaultTeam, false, userID, true); err == nil {
			if _, err = tx.ExecContext(ctx, querySetDefaultTeam, true, userID, teamID); err == nil {
				tx.Commit()
				log.Printf("database: org %s set as default for user %s", teamID, userID)
			} else {
				log.Printf("database: unable to set org %s as default for user %s. Cause: %v", teamID, userID, err)
			}
		} else {
			log.Printf("database: unable to unset default org for user %s. Cause: %v", userID, err)
		}
		if err != nil {
			tx.Rollback()
		}
	}
	return err
}

const queryCreateNewTeam = "INSERT INTO Teams (TeamID, TeamType, TeamName, TeamLanguageCode, TeamStatusID) VALUES (?, ?, ?, ?, ?, ?)"

// CreateNewTeam creates a new User model in the Users collection
func CreateNewTeam(newTeam *model.Team, newRootUser *model.User) error {
	var err error
	var tx *sql.Tx
	ctx := context.TODO()
	if tx, err = database.Connection.BeginTx(ctx, nil); err == nil {
		if _, err = tx.ExecContext(ctx, queryCreateNewTeam, newTeam.TeamID, newTeam.Type, newTeam.Name, newTeam.Language, newTeam.Status); err == nil {
			if err = createNewUser(tx, ctx, newTeam.TeamID, newRootUser, 2); err == nil {
				err = tx.Commit()
			}
		} else {
			log.Printf("database: unable to create user %s. Cause: %v", newTeam.TeamID, err)
		}
		if err != nil {
			tx.Rollback()
		}
	}
	return err
}
