# Paths

### Project:
{
    project_id: int,
    name: string
}

### Bug:
{
    bug_id: int,
    title: string,
    description: string,
    time_amt: double,
    complexity: double,
    project_id: int
}

Routes: 

### /projects
*GET:* Gets all projects 
    Result -  
    [
        {
            Project
        },
        ...
    ]
*POST:* Adds a project into the database
    Data - 
    {
        Project
    }

### /projects/{id}
*GET:* Gets the project at the given {id}
    Result -
    {
        Project
    }
*PUT:* Updates the project at the given {id} with the new data
    Data - 
    {
        Project
    }
*DELETE:* Deletes the project with the given {id}

### /projects/{id}/bugs
*GET:* Gets the bugs from the project with the given {id}
    Result - 
    [
        {
            Bug
        },
        ...
    ]

### /bugs
*GET:* Gets all bugs
    Result - 
    [
        {
            Bug
        },
        ...
    ]
*POST:* Adds a new bug to the databse
    Data -
    {
        Bug
    }

### /bugs/{id}
*GET:* Gets bug with the given {id}
    Result - 
    {
        Bug
    }
*PUT:* Replaces the bug with the given {id} with the data passed in
    Data -
    {
        Bug
    }
*DELETE:* Deletes the bug with the given {id}

