/// <reference path="../pb_data/types.d.ts" />

// Migration: Create events collection for PocketBase v0.23+
migrate((app) => {
    const collection = new Collection({
        "name": "events",
        "type": "base",
        "fields": [
            {
                "name": "title",
                "type": "text",
                "required": true
            },
            {
                "name": "description",
                "type": "text",
                "required": false
            },
            {
                "name": "date_start",
                "type": "date",
                "required": true
            },
            {
                "name": "date_end",
                "type": "date",
                "required": true
            },
            {
                "name": "location",
                "type": "text",
                "required": false
            },
            {
                "name": "url",
                "type": "url",
                "required": true
            },
            {
                "name": "image_url",
                "type": "url",
                "required": false
            },
            {
                "name": "source_name",
                "type": "text",
                "required": true
            },
            {
                "name": "source_id",
                "type": "text",
                "required": true
            },
            {
                "name": "topics",
                "type": "json",
                "required": false
            },
            {
                "name": "category",
                "type": "text",
                "required": true
            },
            {
                "name": "is_new",
                "type": "bool",
                "required": false
            }
        ],
        "listRule": "",
        "viewRule": "",
        "createRule": null,
        "updateRule": null,
        "deleteRule": null
    });

    return app.save(collection);
}, (app) => {
    const collection = app.findCollectionByNameOrId("events");
    return app.delete(collection);
})
