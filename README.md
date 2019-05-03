poddle
======

A music-player PWA written in HTML, (vanilla) JavaScript and CSS, with an
optional backend in Go which supports faster searching, caching, persistence,
and other good stuff.

The point of this project is two-fold: first, podcasts episodes tend to be
distributed as MP3, an ancient digital format which has been superceded
completely by, for example, Opus; and second, to prove that one doesn't need a
"native" app to have a useful app.

You can transcode any podcast episode to opus (using something like the
phenomenal [ffmpeg][]), and get a file that is about 80% smaller than the
equivalent MP3, which is great for mobile devices. Furthermore, your browser
is almost certainly capable of downloading these files and storing them in
their local cache, along with the source of the page, providing a seamless
offline experience.

I got a bit carried away, and ended up making the server quite configurable
(since transcoding on the fly can take some time, the server's timeouts need
to be somewhat flexible), and adding some persistence so that users can, if
they wish, create an account, and store their subscriptions, which will be
kept up to date in a database, mitigating the need to call the feed URL over
and over again.

It might also be cool if a user who has signed in with an SSO provider which
provides cloud storage could _use_ that cloud storage for their podcasts.
Also, since transcoding, while fairly quick, is a multi-minute operation (in
general), it's probably best to persist the files somewhere so they don't need
to be transcoded for every user.

As well as the server (`poddle server`), you can use the executable to
maintain your own local podcast library - see the online help for details, but
the short story is that you specify a database (sqlite3 is fine), `poddle
search` for podcasts to get their feed URLs, `poddle subscribe` to podcasts,
and add a `poddle cron` command to your crontab to automagically check for new
episodes, download them, convert them, store them, and deliver them to, say,
your Telegram account, Dropbox, or email.

[ffmpeg]: https://ffmpeg.org
