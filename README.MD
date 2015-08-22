# Google Contacts

Utility that will extract all contacts from Gmail emails. It takes any emails from To, Cc and saves them into a database. Each time the utility is executed it will scan only new email (email the utility has not seen yet) and add the contact as a new contact if it's not already in the DB

### Parsed fields from email:
The utility parses below fields from email:

- To - to email address
- From - from
- Cc

If the email already exists in the DB, the utility will ignore it and not add or update it. Otherwise it will add it.


### DB Fields:
The utility saves below DB columns:

Contacts
- email (from To, From, Cc)
- first name (from To, From, Cc)
- last name (from To, From, Cc)

Settings
- lastMessageId - internalDate field in Gmail. The internal message creation timestamp (epoch ms), which determines ordering in the inbox. For normal SMTP-received email, this represents the time the message was originally accepted by Google, which is more reliable than the Date header. However, for API-migrated mail, it can be configured by client to be based on the Date header. This is used to we only check messages after this date and don't constantly check the entire mailbox.




## Getting Started

Run the GmailContacts utility.


    $ ./GmailContacts


Once it starts, the utility will first check for security credentials. Specifically, it will ask open a specific link in a browser. The page will ask the user to allow the utility to access Gmail. Once approved the user needs to paste the token back into the utility.


    #########################
    #     Authorization     #
    #########################

    Open below link in your browser. Then Allow the client to access your mailbox,get the token and paste it below

    https://accounts.google.com/o/oauth2/auth?access_type=offline&client_id=4523535455-jdfsfsfdgg7dfadfafdfsdflltghgdk.apps.googleusercontent.com&redirect_uri=urn%3Aieaf%3Afg%3Aoauth%3A2.0%3Aoob&response_type=code&scope=https%3A%2F%2Fwww.googleapis.com%2Fauth%2Fgmail.readonly&state=state-token

    Paste the token here:4/CdflaksdjhfaldkfhaiC9mWGdfadhlfhaslfasdkdasd


Credits
-------
VCImage based on on TJViewController created by Travis Jeffery on 2012-08-06.


Contact
-------
Follow Vladimir Collak on Twitter @vcollak


License
-------

Copyright (c) 2013 Vladimir Collak

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
the documentation of any redistributions of the template files themselves
(but not in projects built using the templates).

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
