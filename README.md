## CrypTag

[![Join the chat at https://gitter.im/elimisteve/cryptag](https://badges.gitter.im/elimisteve/cryptag.svg)](https://gitter.im/elimisteve/cryptag?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)

**_Encrypted, taggable, searchable cloud storage._**

CrypTag was announced at DEF CON 23 in August of 2015.  Presentation
slides:
<https://www.slideshare.net/elimisteve/cryptag-building-encrypted-taggable-searchable-zeroknowledge-systems-59707471>.


### What is CrypTag?

CrypTag is an idea, a spec, an API, and a piece of software that makes
it easy to build a zero-knowledge system, which means that the server
holding user data doesn't know what it is (since it's encrypted).

It is meant as a primitive to be used to build more sophisticated
systems that would rather not re-implement the pieces necessary to
build a zero-knowledge system, but several useful command line
applications have been built with it so far, namely _cput_ (for
encrypting/saving files), _cget_ (for fetching/decrypting files), and
_cpass_ (CryptPass, a password manager).

To use a command line password manager, CryptPass, see "Getting
Started with CryptPass", below.


### How is it searchable _and_ encrypted?

It's not _fully_ searchable; you can query by tag.  See slide 7 of the
presentation from DEF CON 23:
<https://www.slideshare.net/elimisteve/cryptag-building-encrypted-taggable-searchable-zeroknowledge-systems-59707471/8>


### Then the server stores the tags in plaintext?

Nope!  The client stores mapping between tags ("snowden") and a random
hex string ("b6a27d9"), and the server only ever sees the random
strings.

(The client also encrypts these mappings and stores them to the
server, too.)


### Use Cases (what CrypTag is good at) + Syncing via Dropbox

I personally have virtually all data I want shared between my laptops
in one Dropbox folder that CrypTag-based programs add (encrypted) data
to and grab it from.

I've been using cpass to store and fetch...

1. Passwords (**cpass @elimisteve**)
2. Credit card numbers (**cpass visa digits**)
3. Quotes (**cpass nietzsche quote**)
4. Bookmarks, tagged like on Pinboard or Delicious (**cpass url snowden**)
5. Command line commands -- cross-machine shell history! (**cpass install docker**)
6. GitHub recovery codes (**cpass github recoverycode**)

For more on getting started, including how to safely and securely
share passwords with others via a shared Dropbox folder, check out
this section of the README:
https://github.com/cryptag/cryptag#getting-started-with-cryptpass

It's still early days for CrypTag and CryptPass, so don't trust your
life with cpass.  Eventually I will have the code professionally
audited for security flaws.


## Getting Started with CryptPass

The current focus for CrypTag is creating a password manager out of it
called _CryptPass_.  CryptPass exists as a command line tool (`cpass`)
you can use to store and retrieve passwords.  Unencrypted passwords
_never_ touch disk; they are stored encrypted, read into memory, then
printed to your terminal for you to use, with the first one found
added to your clipboard.

### TL;DR version

Install + config:

```
go get github.com/cryptag/cryptag/cmd/cpass
cpass
```

Create passwords, fetch them by tag, or delete them:

```
cpass create mytwitterp4ssw0rd twitter @myusername tag3 tag4
cpass @myusername
cpass delete twitter
```

Keep reading for more advanced options, including password sharing via
shared Dropbox folders.


### Installing cpass

Install Go ([instructions](https://golang.org/doc/install)), then run

    go get github.com/cryptag/cryptag/cmd/cpass

That's it!  Now run

    cpass

`cpass` will generate a new encryption key to store your passwords
with, as well as create the directories it will use to store your
data, all in `~/.cryptag` (by default).


### Using cpass

**Create** a new password and associated tags with commands like:

    cpass create mycr4zyemailp4ssw0rd gmail email elimisteve@gmail

    cpass create mytwitterp4ssw0rd twitter @elimisteve

    cpass create mycr4zyAWSp4ssw0rd4myj0b work aws


**Fetch** the password you're looking for and see them printed to the
screen with commands like:

    cpass gmail

    cpass @elimisteve

    cpass aws work

For convenience, `cpass` adds the first password found to your
clipboard so you can paste it into whichever program you're using.

To **view all** your passwords, run

    cpass all

And finally, to **delete** all passwords with certain tags, run

    cpass delete aws

To only delete one specific password, not all passwords with a generic
tag (e.g., "email") that you may have used to tag multiple passwords,
use the password's tag of the form `id:...`, which is auto-generated
and guaranteed to be unique:

    cpass delete id:a91d46c7-45bb-48e4-43d1-642196df15b2


### Multiple Storage Backends

Maybe you want to store your personal data in `~/.cryptag` but have
passwords you share with colleagues at `~/Dropbox/cryptag_work`, for
example.  (With `cpass` this is secure because plaintext, unencrypted
passwords will never touch `~/Dropbox/cryptag_work`.)

Create a new backend with the desired name (e.g., "work") by running

    BACKEND=work cpass

You can then change the value of `DataPath` at the end of
`~/.cryptag/backends/work.json` to wherever you want your work
passwords stored (e.g., `/home/MYUSERNAME/Dropbox/cryptag_work`),
being sure to use the absolute path.

(Dropbox note: CrypTag-based programs generally, and `cpass`
specifically, store each piece of data (e.g., each password and each
tag) in a separate file, so it _is_ safe for multiple people to create
passwords simultaneously and save them to a shared Dropbox folder,
unlike with KeePass, KeePassX, 1Password, and some other password
managers.)

Now you can save shared work passwords with the same commands as
before, except with the `BACKEND` environment variable
set:

    BACKEND=work cpass create mycr4zyAWSp4ssw0rd4myj0b work aws

Now you should share `~/.cryptag/backends/work.json` with your
colleagues -- or at least the encryption key -- so that you can
decrypt passwords saved by each other.


### More Convenient Multiple Storage Backends

See [this issue](https://github.com/cryptag/cryptag/issues/18) for
discussion on how to make storing data in multiple places much better!
I would love your input.

I currently do this:

    echo 'BACKEND=work cpass "$@"' > ~/bin/work
    chmod +x ~/bin/work

so I can simply do

    work create mycr4zyAWSp4ssw0rd4myj0b aws

to create work passwords, or

    work aws

to fetch them.


### Non-use Cases (what CrypTag is not good at)

Anything that requires rapid changes being made to data by multiple
users, such as:

* Real-time collaborative document editing

  * Real-time spreadsheet editing should work OK (as long as you're OK
    with "last write wins" to a cell), since each cell can be its own
    Row that can be changed concurrently with other Rows

Any data that is "write once, read many times" is ideal for CrypTag.


## Future Plans

I have big plans for CryptPass and other CrypTag-based software to
help make the world's data -- passwords, everything stored "in the
cloud", file backups, bookmarks and so on -- more secure.

I believe that you should be able access your data from any of your
devices, and just grab what you need, exactly like you can from the
Dropbox mobile app.  Or if you don't mind storing all your data on
your computer, being able to use Dropbox (or anything similar) to sync
all your data between all your devices _without having to trust the
company storing your data for you_ is also deeply important; we should
all benefit from the convenience of cloud storage without giving up
any privacy whatsoever.

So whatever feedback you may have, please please send it my way!  Yes,
there will be a **graphical version of CryptPass** usable on Windows,
Mac OS X, and Linux desktops.  Eventually I'd like to have mobile
apps, too, of course.

I am open to **all** questions, comments, suggestions, insults, and
whatever else you've got.


## Geeky Feedback Requested

The graphical version of CryptPass (that uses Electron + React, that
then talks to a local CrypTag daemon) once the command line version,
cpass, is better, and once more complex [storage
questions](https://github.com/cryptag/cryptag/issues/18) are
answered, which I'd appreciate feedback on from those of you who may
want to store different kinds of data in different places (e.g., all
passwords in a local directory, all work passwords in a shared Dropbox
folder, and all backups in S3).

I'd love to create mobile versions of CryptPass, probably starting
with Ubuntu Phone, because I can write it all in Go :-), and now that
both Android and iOS apps can call into Go code using some [new
awesome mobile
shit](https://godoc.org/golang.org/x/mobile/cmd/gomobile), it
shouldn't be necessary to port the core CrypTag logic to another
language.

Thank you!  Here's to a more privacy-friendly future for all!


# Cryptography Notice

This distribution includes cryptographic software. The country in which you currently reside may have restrictions on the import, possession, use, and/or re-export to another country, of encryption software.
BEFORE using any encryption software, please check your country's laws, regulations and policies concerning the import, possession, or use, and re-export of encryption software, to see if this is permitted.
See <http://www.wassenaar.org/> for more information.

The U.S. Government Department of Commerce, Bureau of Industry and Security (BIS), has classified this software as Export Commodity Control Number (ECCN) 5D002.C.1, which includes information security software using or performing cryptographic functions with asymmetric algorithms.
The form and manner of this distribution makes it eligible for export under the License Exception ENC Technology Software Unrestricted (TSU) exception (see the BIS Export Administration Regulations, Section 740.13) for both object code and source code.
