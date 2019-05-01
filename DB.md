The database:

podcasts
id (uuid)
url (url)
title (string)
description (text)
image (image)
created_at (datetime)
updated_at (datetime)
published_at (datetime)

episodes
id (uuid)
podcast_id (uuid)
title (string)
description (text)
image (image)
created_at (datetime)
updated_at (datetime)
published_at (datetime)

versions
id (uuid)
url (url)
podcast_id (uuid)
episode_id (uuid)
type (string)
created_at (datetime)
updated_at (datetime)
progress (float)

users
id (uuid)
subscriptions ([]uuid)

identities
id (uuid)
url (url)
provider (provider)
access_token (string)
refresh_token (string)
data (json)
