module default {
    type Todo {
        required property item -> str;
        property done_at -> datetime;
        required property updated_at -> datetime;
        required property created_at -> datetime;
    }
}
