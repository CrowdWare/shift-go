# shift-go
This is the business logic of the Shift Android app.  
We are writing this in Go because of the fact that every developer is able to recompile a Kotlin based app and is then able to see passwords, API keys and the like.  
Not so easy with Go compiled code.

My first idea was to keep this library close source, to make it a bit harder for hackers, but something reminds me, why I decided to make everything open source.  
Thrust.  
I want my users to trust our apps.  
- No hidden secrets
- No spying
- Let other developers find mistakes

So today I open up the source code.

# run
In order to run you should rename crypto_vars.go.sample to cyrpto_vars.go and edit the values, like you will find in release.sh.