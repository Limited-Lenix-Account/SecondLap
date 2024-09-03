# What is this?
This is a raffle bot I originally made for Motorsport's KTM bike raffle. It includes an sql database to keep track of entries, proxy support, and uses 2cap to solve captchas. I made this just for fun not really thinking I'd win anything (which i didn't :( ), but thought it'd give me the opportunity to learn something new. This is my first time using SQL so if I didn't follow any best practices, my bad. I've been making small scripts for sites with ephemeral purposes for a while but this is the first one I thought was decent enough to make public after the raffle closed.

# What do I Need to Run It?
The only files you'll have to make yourself are the `proxies.txt` and `emails.txt` files in the `data` folder. You'll have to propagate with their respective information. A database will created that adds the email to the SQL database that will be read from when you start your tasks.

Additionally, you'll need a 2cap key to solve captchas for each task that you run.

# License
[MIT](https://choosealicense.com/licenses/mit/)
