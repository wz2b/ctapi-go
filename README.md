# Aveva Plant SCADA (formerly CiTect SCADA) CTAPI for Go

This project binds to the 
[AVEVA Plant SCADA](https://www.aveva.com/en/products/plant-scada/)
runtime from go.  It provides bindings
for all of the functions, and a reasonable facade
for the ones that are most commonly used.

# How to use this library

## Initialization

To initialize the library, first get an instance of CtApi
by calling Init().  Pass in the name of the DLL, and the
path to the directory containing it.  Note that I have done
all my development with the 64-bit version - it may be
possible to do this with a 32-bit version as well if you
build for that.

```go
var dllPath = "C:\\Program Files (x86)\\AVEVA Plant SCADA\\Bin\\Bin (x64)"
api, err := ctapi.Init("CtApi.dll", dllPath)
if err != nil {
    panic(err)
}
```

The initialization will open the DLL and find the exported
functions.  At this point, you can open a connection to
Plant SCADA:

```
handle, err := api.CtOpen()
if err != nil {
    panic(err)
} else {
    println(handle)
}
defer api.CtClose(handle)
```

## Find
The Find API is a little odd to use, so I have wrapped it in
something as close to a _generator_ as go allows

```go
const tableName = "Tag"
const propertyName = "TAG"

for hObj := range api.FindAll(handle, tableName) {
  tag, e := api.GetStringProperty(hObj, propertyName)
  if e != nil {
      fmt.Printf("There was an error getting the string property: %v\n", e)
  } else {
      fmt.Printf("%v \n", tag)
  }
}	
	
```

This enumerates all items from the _Tag_ table, then gets the property named *TAG* and
prints it to the screen.  The *TAG* property is the tag name. You can use `FindAll` to
find other sorts of things including _Alarm_ (alarm tags), _DigAlm_ (digital alarm tags),
_Trend_, and a few other types.  Each table has its own unique set of properties.  To
find out what they are, refer to the [Browse Function Field Reference](https://gcsresource.aveva.com/plantscada/WebHelp/plantscada2023/Content/Cicode/Browse_Function_Field_Reference.html)
documentation (you have to have an AVEVA support account for this link to work).


