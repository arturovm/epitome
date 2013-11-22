# Pond API

## Table of Contents

1. [General Considerations](#general-consideration)
2. [Users](#users)
3. [Authentication](#authentication)
4. [Subscriptions](#subscriptions)
5. [Articles](#articles)
6. [Preferences](#preferences)

## General Considerations

All request paths must be prepended with the server's API root. In the case of the Pond reference implementation, this is `/api`.

Except where noted, each request must include the session token obtained during authentication, in the `X-Session-Token` HTTP header.

When the server encounters an error, it will respond with the apprpriate HTTP status code, and a JSON object containing a single `error` field.

Also keep in mind that some users will not have SSL/TLS available.

## Users

__Note:__ The users API is still incomplete, but all the essential functionality is there (i.e. creating users). In the near future (quite possibly, in the next minor version), methods will be added for deleting and editing users.

### Creating Users

Depending on server preferences, a session token may not be required to perform this action.

#### `POST /users`

#### Form Variables

|  Variable  | Value             | Required | Default  | Explanation |
|------------|-------------------|----------|----------|-------------|
| `username` | string            | __Yes__  | _none_   | The user's chosen username _as–is_ (i.e. no need to convert it to lowercase). |
| `password` | string            | __Yes__  | _none_   | The password hash, as per the procedure explained below. |
| `role`     | `normal`, `admin` | No       | `normal` | The user's role, where the role is equal or greater than the minimum role required for account creation, specified in the server's preferences. |

#### Response

| Code  | Body                                           | Explanation |
|-------|------------------------------------------------|-------------|
| `201` | _none_                                         | The user was created successfully. |
| `409` | [A JSON error object](#general-considerations) | A user with this username already exists. |

#### Generating a Password Hash

To generate the password hash that must be sent to the server, concatenate the user's username converted to lowercase with the user's chosen password placing a colon (`:`) in between, and calculate the MD5 hash of the result.

```
passwordHash = MD5(lowercase(username) + ':' + password)
```

#### Notes

Currently, there is a vulnerability in the design of the first–time setup that allows anybody form the Internet to create a new user with an `admin` role. This will be fixed in the next patch version (`0.1.1`).

## Authentication

### Logging In (i.e. Creating a Session)

A session token is not required to perform this action.

#### `POST /auth/sessions`

#### Form Variables

| Variable   | Value  | Required | Default | Explanation |
|------------|--------|----------|---------|-------------|
| `username` | string | __Yes__  | _none_  | The user's chosen username. This operation is case–insensitive. |
| `password` | string | __Yes__  | _none_  | The password hash, as per the procedure explained in [Generating a Password Hash](#generating-a-password-hash). |

#### Response

| Code  | Body                                                            | Explanation |
|-------|-----------------------------------------------------------------|-------------|
| `201` | A JSON object containing a single field, called `session_token` | The user was logged in successfully. |
| `400` | A JSON error object                                             | Username or passsword were not included in the request. |
| `401` | A JSON error object                                             | Wrong password. |
| `404` | A JSON error object                                             | A user with the provided username doesn't exist. |

### Logging Out (i.e. Deleting a Session)

#### `DELETE /auth/sessions/:session-token`

#### Path Components

| Variable         | Explanation |
|------------------|-------------|
| `:session-token` |  The session token itself. |

#### Response

| Code  | Body                | Explanation |
|-------|---------------------|-------------|
| `200` | _none_              | The user was logged out successfully. |
| `400` | A JSON error object | The session token was not provided. |

## Subscriptions

To avoid data redundancy, Pond keeps two tables: one which stores server–wide subscriptions, and another which stores per–user _references_ to the server–wide table. Whenever a user subscribes to a feed, the server first checks if it already exists in the server–wide table; if it doesn't, it adds it, and then adds a user–specific reference to it; if it does, it just adds the reference.

You can delete server–wide subscriptions by setting the `global` query parameter to `true`, when [unsubscribing from a feed](#unsubscribing-from-a-feed). Be _very_ careful when doing this; it affects _all_ users subscribed to that feed.

### Subscribing to a Feed

When adding a subscription, if the URL is not a feed itself but rather a website, the server attempts to perform automatic discovery of the RSS or Atom feed, via the `<link>` HTML tags.

#### `POST /subscriptions`

#### Form Variables

| Variable | Value  | Required | Default | Explanation |
|----------|--------|----------|---------|-------------|
| `url`    | string | __Yes__  | _none_  | This could be an Atom or RSS feed URL, or the URL of a website. |

#### Response

| Code  | Body                | Explanation |
|-------|---------------------|-------------|
| `201` | _none_              | The subscription was added correctly. |
| `400` | A JSON error object | The session token was not provided or the url variable was not set. |
| `404` | A JSON error object | The feed URL could not be discovered. |
| `409` | A JSON error object | The user is already subscribed to that feed. |

### Getting all Subscriptions

#### `GET /subscriptions`

#### Form Variables

| Variable | Value   | Required | Default | Explanation |
|----------|---------|----------|---------|-------------|
| `global` | boolean | No       | `false` | Whether to return global server subscriptions. Admin privileges are required to perform this action. |

#### Response

| Code  | Body                                                                            | Explanation |
|-------|---------------------------------------------------------------------------------|-------------|
| `200` | A possibly empty JSON array of [Subscription objects](#the-subscription-object) | Returns a list of the user's subscriptions. |
| `400` | A JSON error object                                                             | The session token was not provided. |
| `401` | A JSON error object                                                             | If `global` was set to `true`, this might mean the user doesn't have enough privileges to view global subsciptions. |

### Unsubscribing from a Feed

#### `DELETE /subscriptions/:id`

#### Path Components

| Variable | Explanation |
|----------|-------------|
| `:id`    | ID of the subscription from which you wish to unsubscribe. |

#### Form Variables

| Variable | Value   | Required | Default | Explanation |
|----------|---------|----------|---------|-------------|
| `global` | boolean | No       | `false` | Whether to unsubscribe from the feed globally. Admin privileges are required to perform this action. |

#### Response

| Code  | Body                | Explanation |
|-------|---------------------|-------------|
| `200` | _none_              | The user was unsubscribed successfully from the feed. |
| `400` | A JSON error object | The session token was not provided. |
| `404` | A JSON error object | The subscription doesn't exist. |

### The Subscription Object

| Field Name | Value   |
|------------|---------|
| `id`       | integer |
| `url`      | string  |
| `name`     | string  |

## Articles

### Fetching Articles

#### `GET /subscriptions/articles`

#### Query Parameters

| Variable    | Value                   | Required | Default | Explanation |
|-------------|-------------------------|----------|---------|-------------|
| `status`    | `all`, `unread`, `read` | No       | `all`   | Whether to return all articles, unread articles only, or read articles only. |
| `limit`     | integer                 | No       | `100`   | How many articles to return. |
| `order`     | `desc`, `asc`           | No       | `desc`  | Whether to return newer articles first (`desc`) or older articles first (`asc`). |
| `since_id`  | integer                 | No       | _none_  | If set, returns all articles published after the article with the specified ID. |
| `before_id` | integer                 | No       | _none_  | If set, returns all articles published before the article with the specified ID. |

#### Response

| Code  | Body                                                                  | Explanation |
|-------|-----------------------------------------------------------------------|-------------|
| `200` | A possibly empty JSON array of [Article objects](#the-article-object) | Returns the user's articles, according to the parameters of the query. |
| `400` | A JSON error object                                                   | The session token was not provided. |

### Fetching Articles from a Specific Feed

#### `GET /subscriptions/:id/articles`

#### Path Components

| Variable | Explanation |
|----------|-------------|
| `:id`    | The ID of the feed from which to fetch articles. |

#### Query Parameters

| Variable    | Value                   | Required | Default | Explanation |
|-------------|-------------------------|----------|---------|-------------|
| `status`    | `all`, `unread`, `read` | No       | `all`   | Whether to return all articles, unread articles only, or read articles only. |
| `limit`     | integer                 | No       | `100`   | How many articles to return. |
| `order`     | `desc`, `asc`           | No       | `desc`  | Whether to return newer articles first (`desc`) or older articles first (`asc`). |
| `since_id`  | integer                 | No       | _none_  | If set, returns all articles published after the article with the specified ID. |
| `before_id` | integer                 | No       | _none_  | If set, returns all articles published before the article with the specified ID. |

#### Response

| Code  | Body                                                                  | Explanation |
|-------|-----------------------------------------------------------------------|-------------|
| `200` | A possibly empty JSON array of [Article objects](#the-article-object) | Returns the user's articles, according to the parameters of the query. |
| `400` | A JSON error object                                                   | The session token was not provided. |
| `404` | A JSON error object                                                   | The subscription doesn't exist. |

### Marking an Article as Read or Unread

#### `PUT /subscriptions/:subscriptionid/articles/:articleid`

#### Path Components

| Variable          | Explanation |
|-------------------|-------------|
| `:subscriptionid` | The ID of the feed to which the article belongs. |
| `:articleid`      | The ID of the article itself. |

#### Query Parameters

| Variable | Value            | Required | Default | Explanation |
|----------|------------------|----------|---------|-------------|
| `status` | `read`, `unread` | __Yes__  | _none_  | The status to which you wish to change the article. |

#### Response

| Code  | Body   | Explanation |
|-------|--------|-------------|
| `200` | _none_ | The article was successfully marked as read or unread. |
| `400` | A JSON error object | The session token was not provided or an invalid value was passed to `status`. |
| `404` | A JSON error object | The article or subscription doesn't exist. |

### The Article Object

| Field Name        | Value   |
|-------------------|---------|
| `id`              | integer |
| `subscription_id` | integer |
| `url`             | string  |
| `title`           | string  |
| `author`          | string  |
| `published_at`    | ISO 8601 UTC+0 formatted string |
| `summary`         | object  |
| `summary.type`    | `xhtml`, `html`, `text` |
| `summary.content` | string  |
| `body`            | object  |
| `body.type`       | `xhtml`, `html`, `text` |
| `body.content`    | string  |
| `read`            | boolean |

## Preferences

Admin privileges are required to view or edit server preferences.

### Getting Current Preferences

#### `GET /preferences`

#### Response

| Code  | Body                                            | Explanation |
|-------|-------------------------------------------------|-------------|
| `200` | A [Preferences object](#the-preferences-object) | Returns the current server preferences. |
| `400` | A JSON error object                             | The session token was not provided. |

### Changing Preferences

#### `PUT /preferences`

#### Request Body

A [Preferences object](#the-preferences-object) containing at least one field.

#### Response

| Code  | Body                | Explanation |
|-------|---------------------|-------------|
| `200` | _none_              | Server preferences were successfuly changed. |
| `400` | A JSON error object | The session token was not provided, or no fields were sent within the JSON object, or an invalid value was sent for one of the fields. |

### The Preferences Object

| Field Name             | Value   | Explanation |
|------------------------|---------|-------------|
| `refresh_rate`         | string  | A string representing an interval, in the format of `(([0-9])+h)?(([0-9])+m)?(([0-9])+s)?`, where `s` is the "seconds" identifier, `m` is the "minutes" identifier and `h` is the "hours" identifier. |
| `new_user_permissions` | integer | An integer representing the minimum role required to create new users, in the range of [0, 2], where 0 is `admin`, 1 is `normal` and 2 is `public` |
