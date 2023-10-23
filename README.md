# Aveva Plant SCADA (formerly CiTect SCADA) CTAPI for Go

This project binds to the
[AVEVA Plant SCADA](https://www.aveva.com/en/products/plant-scada/)
runtime from go. It provides bindings
for all of the functions, and a reasonable facade
for the ones that are most commonly used.

# License

I am releasing this under an Apache 2.0 license, but if you need it to
be something else please contact me. Consider this "postcard-ware" in
that I'd appreciate a note if you find this useful.

# Why did I build this

SCADA software tends to be frustrating if you're a software engineer.
I think the reason for this is that these suites were designed with
those who have a strong controls background but not necessarily depth
of knowledge in software development.  The result is often something
that's easy for the target audience to use, but limiting for those
who want to extend functionality beyond what the vendor had in mind.

This library exists because the default language for customizing
things in Plant SCADA is a basic-like language named CiCode, and it
is pretty limited. It can do a lot but can't handle things like
REST calls, JSON/XML serialization, and having a sane way to open sockets
to other services. You can send e-mail from CiCode but only through
MAPI (pretty obsolete). In an attempt to bring Plant SCADA into the
current decade, and provide greater integration opportunities, I wrote
this library. My main intention here is:

* Integrating with some external IIoT frameworks like Amazon SiteWise
* Linking to a more robust paging solution that supports duty schedules,
  resolution tracking, and Android/IOS Push API. Some people like
  [PagerDuty](https://www.pagerduty.com).  [Squadcast](https://www.squadcast.com/)
  is also a good option that has easy integrations and a lot of features.
  Grafana [OnCall](https://grafana.com/products/cloud/oncall) is newer but
  it looks like they have the right idea.
* Storing data in a more reasonable time-series database like
  TimescaleDB or Influx, so that I can have much more flexibility
  visualizing data


Of course, all of this could have been done in Visual Studio (but why?)

# How to use this library

## Initialization

To initialize the library, first get an instance of CtApi
by calling Init(). Pass in the name of the DLL, and the
path to the directory containing it. Note that I have done
all my development with the 64-bit version - it may be
possible to do this with a 32-bit version as well if you
build for that - I haven't tried, though, and I'm a little
worried about how it will handle UTF16 strings. If you try
it and it works (or doesn't) please drop me a note and let me know.

```go
var dllPath = "C:\\Program Files (x86)\\AVEVA Plant SCADA\\Bin\\Bin (x64)"
api, err := ctapi.Init("CtApi.dll", dllPath)
if err != nil {
panic(err)
}
```

The initialization will open the DLL and find the exported
functions. At this point, you can open a connection to
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
prints it to the screen. The *TAG* property is the tag name. You can use `FindAll` to
find other sorts of things including _Alarm_ (alarm tags), _DigAlm_ (digital alarm tags),
_Trend_, and a few other types. Each table has its own unique set of properties. To
find out what they are, refer to
the [Browse Function Field Reference](https://gcsresource.aveva.com/plantscada/WebHelp/plantscada2023/Content/Cicode/Browse_Function_Field_Reference.html)
documentation (you have to have an AVEVA support account for this link to work).


