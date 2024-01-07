# Objects

### Project:
<pre>
{
    projectID: int,
    name: string
}
</pre>

### Bug:
<pre>
{
    bugID: int,
    title: string,
    description: string,
    timeAmt: double,
    complexity: double,
    projectID: int
}
</pre>

# Endpoints: 

### /projects
**GET:** Gets all projects 
<pre>
    Result -  
    [
        {
            Project
        },
        ...
    ]
</pre>

**POST:** Adds a project into the database
<pre>
    Data - 
    {
        Project
    }
</pre>

### /projects/{id}
**GET:** Gets the project at the given {id}
<pre>
    Result -
    {
        Project
    }
</pre>

**PUT:** Updates the project at the given {id} with the new data
<pre>
    Data - 
    {
        Project
    }
</pre>
**DELETE:** Deletes the project with the given {id}

### /projects/{id}/bugs
**GET:** Gets the bugs from the project with the given {id}
<pre>
    Result - 
    [
        {
            Bug
        },
        ...
    ]
</pre>

### /bugs
**GET:** Gets all bugs
<pre>
    Result - 
    [
        {
            Bug
        },
        ...
    ]
</pre>
**POST:** Adds a new bug to the databse
<pre>
    Data -
    {
        Bug
    }

    OR for batch insert

    [
        {
            Bug
        },
        {
            Bug
        }, ...
    ]
</pre>

### /bugs/{id}
**GET:** Gets bug with the given {id}
<pre>
    Result - 
    {
        Bug
    }
</pre>
**PUT:** Replaces the bug with the given {id} with the data passed in
<pre>
    Data -
    {
        Bug
    }
</pre>
**DELETE:** Deletes the bug with the given {id}