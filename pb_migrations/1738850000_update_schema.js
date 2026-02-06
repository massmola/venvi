/// <reference path="../pb_data/types.d.ts" />

migrate((app) => {
    // 1. Update 'events' collection
    const events = app.findCollectionByNameOrId("events");

    events.fields.add(new NumberField({
        "name": "latitude",
        "required": false
    }));
    events.fields.add(new NumberField({
        "name": "longitude",
        "required": false
    }));

    app.save(events);

    // 2. Update 'users' collection
    const users = app.findCollectionByNameOrId("users");

    users.fields.add(new TextField({
        "name": "name", // standard OAuth mapping usually maps to 'name'
        "required": false
    }));
    users.fields.add(new NumberField({
        "name": "latitude",
        "required": false
    }));
    users.fields.add(new NumberField({
        "name": "longitude",
        "required": false
    }));

    app.save(users);

}, (app) => {
    // Revert 'events' additions
    const events = app.findCollectionByNameOrId("events");
    events.fields.removeByName("latitude");
    events.fields.removeByName("longitude");
    app.save(events);

    // Revert 'users' additions
    const users = app.findCollectionByNameOrId("users");
    users.fields.removeByName("name");
    users.fields.removeByName("latitude");
    users.fields.removeByName("longitude");
    app.save(users);
})
