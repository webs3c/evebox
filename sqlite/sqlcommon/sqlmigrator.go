// +build cgo

/* Copyright (c) 2016 Jason Ish
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions
 * are met:
 *
 * 1. Redistributions of source code must retain the above copyright
 *    notice, this list of conditions and the following disclaimer.
 * 2. Redistributions in binary form must reproduce the above copyright
 *    notice, this list of conditions and the following disclaimer in the
 *    documentation and/or other materials provided with the distribution.
 *
 * THIS SOFTWARE IS PROVIDED ``AS IS'' AND ANY EXPRESS OR IMPLIED
 * WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF
 * MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
 * DISCLAIMED. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR ANY DIRECT,
 * INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
 * (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
 * SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION)
 * HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT,
 * STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING
 * IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
 * POSSIBILITY OF SUCH DAMAGE.
 */

package sqlcommon

import (
	"database/sql"
	"fmt"
	"github.com/jasonish/evebox/log"
	"github.com/jasonish/evebox/resources"
	"os"
	"path"
)

type SqlMigrator struct {
	db        *sql.DB
	directory string
}

func NewSqlMigrator(db *sql.DB, directory string) *SqlMigrator {
	return &SqlMigrator{
		db:        db,
		directory: directory,
	}
}

func (m *SqlMigrator) Migrate() error {

	var currentVersion int
	nextVersion := 0

	rows, err := m.db.Query("select max(version) from schema")
	if err == nil {
		if rows.Next() {
			if err := rows.Scan(&currentVersion); err != nil {
				return err
			}
			nextVersion = currentVersion + 1
		}
		rows.Close()
		log.Debug("Current database schema version: %d", currentVersion)
	} else {
		log.Debug("Initializing database.")
	}

	for {
		scriptName := path.Join(m.directory,
			fmt.Sprintf("V%d.sql", nextVersion))
		script, err := resources.AssetString(scriptName)
		if err != nil {
			break
		}

		log.Info("Updating database to version %d.", nextVersion)

		tx, err := m.db.Begin()
		if err != nil {
			log.Error("Failed to start transaction: %v", err)
			return err
		}

		_, err = tx.Exec(script)
		if err != nil {
			log.Error("Failed to execute script: %v", err)
		}

		err = m.setVersion(tx, nextVersion)
		if err != nil {
			log.Error("Failed to update schema version: %v", err)
			return err
		}

		err = tx.Commit()
		if err != nil {
			log.Error("Fail to commit transaction: %v", err)
			return err
		}

		nextVersion++
	}

	return nil
}

func (m *SqlMigrator) setVersion(tx *sql.Tx, version int) error {
	_, err := tx.Exec(`insert into schema (version, timestamp)
	                     values ($1, datetime('now'))`, version)
	return err
}

func (m *SqlMigrator) fileExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	return false
}
