// Go-lang password manager
//
// Pedro Pardo, ppardo3
// U. of Illinois, Chicago
// CS 341, Fall 2023
//

package main

import (
	"fmt"
	"os"
)

// Global variables are allowed (and encouraged) for this project.

type Entry struct {
	Site     string
	User     string
	Password string
}

type EntrySlice []Entry

var passwordMap map[string]EntrySlice

const passwordVaultFile = "passwordVault"

//_______________________________________________________________________
// initialize before main()
//_______________________________________________________________________
func init () {
	passwordMap = make(map[string]EntrySlice) //create map for passwords
}

//_______________________________________________________________________
// find the matching entry slice
//_______________________________________________________________________
func findEntrySlice(site string) (EntrySlice, bool) {
	entrySlice, found := passwordMap[site] //finds entry slice for given site
	return entrySlice, found
}

//_______________________________________________________________________
// set the entrySlice for site
//_______________________________________________________________________
func setEntrySlice(site string, entrySlice EntrySlice)  {
	passwordMap[site] = entrySlice //sets entry slice for site
}

//_______________________________________________________________________
// find
//_______________________________________________________________________
func find(user string, entrySlice EntrySlice) (int, bool) {
	for i, entry := range entrySlice { //finds index within the entryslice
		if entry.User == user {
			return i, true
		}
	}
	return -1, false
}

//_______________________________________________________________________
// print the list in columns
//_______________________________________________________________________
func pmList() {
	fmt.Printf("%-20s %-20s %-20s\n", "Site", "User", "Password") 
	fmt.Println("------------------------------------------------------------") //output ststaments for passwords
	for _, entrySlice := range passwordMap { 
		for _, entry := range entrySlice {
			fmt.Printf("%-20s %-20s %-20s\n", entry.Site, entry.User, entry.Password) //prints what is ever store in the map
		}
	}
}

//_______________________________________________________________________
//  add an entry if the site, user is not already found
//_______________________________________________________________________
func pmAdd(site, user, password string) {
	entrySlice, found := findEntrySlice(site)
	if found {
		_, userFound := find(user, entrySlice) //checks for user
		if userFound {
			fmt.Println("add: duplicate entry") //if found output error
			return
		}

		// Add the new entry for the existing site
		entry := Entry{Site: site, User: user, Password: password}
		passwordMap[site] = append(entrySlice, entry)
	} else {
		entry := Entry{Site: site, User: user, Password: password} //add new site
		passwordMap[site] = EntrySlice{entry}
	}
}

//_______________________________________________________________________
// remove by site and user
//_______________________________________________________________________
func pmRemove(site, user string) {
	entrySlice, found := findEntrySlice(site) //searched for website
	if !found {
		fmt.Println("remove: site not found") //error for no website being stored
		return
	}

	index, userFound := find(user, entrySlice) //finds user if first bit is valid
	if userFound {
		passwordMap[site] = append(entrySlice[:index], entrySlice[index+1:]...) //removes from map
	} else {
		fmt.Println("remove: user not found") //print out error
	}
}

//_______________________________________________________________________
// remove the whole site if there is a single user at that site
//_______________________________________________________________________
func pmRemoveSite(site string) {
	entrySlice, found := findEntrySlice(site) //finds entryslice
	if !found {
		fmt.Println("remove: site not found") //error if returned false
		return
	}

	if len(entrySlice) == 1 {
		delete(passwordMap, site) //checks if users in site is more then one
	} else {
		fmt.Println("attempted to remove multiple users") //cant remove cause more then 1
	}

}


//_______________________________________________________________________
// read the passwordVault
//_______________________________________________________________________
func pmRead() {
	file, err := os.Open(passwordVaultFile) //opens file
	if err != nil {
		fmt.Println("Error reading passwordVault:", err) //cant open it
		return
	}
	defer file.Close()

	var entry Entry
	for {
		_, err := fmt.Fscanf(file, "%s %s %s\n", &entry.Site, &entry.User, &entry.Password)
		if err != nil {
			break // End of file
		}

		passwordMap[entry.Site] = append(passwordMap[entry.Site], entry) //update
	}
}

//_______________________________________________________________________
// write the passwordVault
//_______________________________________________________________________
func pmWrite() {
	file, err := os.Create(passwordVaultFile) //create file
	if err != nil {
		fmt.Println("Error writing to passwordVault:", err) //error for password Vault
		return
	}
	defer file.Close()

	for _, entrySlice := range passwordMap {
		for _, entry := range entrySlice {
			fmt.Fprintf(file, "%s %s %s\n", entry.Site, entry.User, entry.Password) //prints file
		}
	}
}



//_______________________________________________________________________
// do forever loop reading the following commands
//    l
//    a s u p
//    r s
//    r s u
//    x
//  where l,a,r,x are list, add, remove, and exit
//  and s,u,p are site, user, and password
//_______________________________________________________________________
func loop() {
	for {
		var command string //string for input
		fmt.Print("Enter command (l/a/r/x): ")
		fmt.Scan(&command) //scans in input

		switch command { //switch case to run valid commands
		case "l":
			pmList() //calls list funciton
		case "a": //takes in new password
			var site, user, password string
			fmt.Scan(&site, &user, &password) //takes input
			pmAdd(site, user, password) //runs add func
		case "r":
			readRemoveCommand() //runs remove helper function
		case "x":
			pmWrite() //write to file and saves
			os.Exit(0)
		default:
			fmt.Println("Invalid command. Please enter l, a, r, or x.") //if input isnt valid run again
		}
	}
}

func readRemoveCommand() { //remove helper functoin
	var site, user string //inputs

	if _, err := fmt.Scanf("%s %s\n", &site, &user); err != nil { //this checks if the input is just 1 string or two 
		pmRemoveSite(site) //if inpout is just 1 then run remove site
	} else {
		pmRemove(site, user) //when input is 2
	}
}

//_______________________________________________________________________
//  let her rip
//_______________________________________________________________________
func main() {
	pmRead() //read file if there
	loop() //run main loop
}
