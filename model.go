package main

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
)

type ID int

type Person struct {
	ID        ID        `db:"id" json:"id,omitempty"`
	FirstName string    `db:"first_name" json:"first_name,omitempty"`
	LastName  string    `db:"last_name" json:"last_name,omitempty"`
	DOB       time.Time `db:"dob" json:"dob,omitempty"`

	Employers []Company `db:"employers" json:"employers,omitempty"`
}

func (p Person) String() string {
	var str string
	str += fmt.Sprintf("[%d] %s %s", p.ID, p.FirstName, p.LastName)
	if len(p.Employers) > 0 {
		str += "\n    Employers:"
	}
	for _, empl := range p.Employers {
		str += fmt.Sprintf("\n   - [%d] %s", empl.ID, empl.Name)
	}

	return str
}

type Company struct {
	ID   ID     `db:"id" json:"id,omitempty"`
	Name string `db:"name" json:"name,omitempty"`

	Employees []Person `db:"employees" json:"employees,omitempty"`
}

func (p Company) String() string {
	var str string
	str += fmt.Sprintf("[%d] %s", p.ID, p.Name)
	if len(p.Employees) > 0 {
		str += "\n    Employees:"
	}
	for _, empl := range p.Employees {
		str += fmt.Sprintf("\n   - [%d] %s %s", empl.ID, empl.FirstName, empl.LastName)
	}

	return str
}

const getPersonSQL = `
	select id, first_name, last_name,
		(
			select array_agg(row(company.id, company.name))
			from company
			join company_person on company.id = company_person.company_id
			where company_person.person_id = person.id
		) as employers
	from person
	where id = $1
`

func GetPerson(ctx context.Context, db DB, id ID) (Person, error) {
	rows, _ := db.Query(ctx, getPersonSQL, id)
	person, err := pgx.CollectOneRow[Person](rows, pgx.RowToStructByNameLax)
	return person, err
}

const getPeopleSQL = `
	select
		id,
		first_name,
		last_name,
		(
			select array_agg(row(company.id, company.name))
			from company
			join company_person on company.id = company_person.company_id
			where company_person.person_id = person.id
		) as employers
	from person
`

func GetPeople(ctx context.Context, db Queryer) ([]Person, error) {
	rows, _ := db.Query(ctx, getPeopleSQL)
	people, err := pgx.CollectRows[Person](rows, pgx.RowToStructByNameLax)
	return people, err
}

const getCompanySQL = `
	select id, name, 
		(
			select array_agg(row(person.id, person.first_name, person.last_name))
			from person
			join company_person on person.id = company_person.person_id
			where company_person.company_id = company.id
		) as employees
	from company
	where id = $1
`

func GetCompany(ctx context.Context, db DB, id ID) (Company, error) {
	rows, _ := db.Query(ctx, getCompanySQL, id)
	company, err := pgx.CollectOneRow[Company](rows, pgx.RowToStructByNameLax)
	return company, err
}

const getCompaniesSQL = `
	select 
		id,
		name,
		(
			select array_agg(row(person.id, person.first_name, person.last_name))
			from person
			join company_person on person.id = company_person.person_id
			where company_person.company_id = company.id
		) as employees
	from company
`

func GetCompanies(ctx context.Context, db Queryer) ([]Company, error) {
	rows, _ := db.Query(ctx, getCompaniesSQL)
	companies, err := pgx.CollectRows[Company](rows, pgx.RowToStructByNameLax)
	return companies, err
}
