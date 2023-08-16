# Objects

### Project:
<pre>
{
    project_id: int,
    name: string
}
</pre>

### Bug:
<pre>
{
    bug_id: int,
    title: string,
    description: string,
    time_amt: double,
    complexity: double,
    project_id: int
}
</pre>

# Routes: 

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