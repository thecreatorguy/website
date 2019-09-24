<!DOCTYPE html>
<html lang="en-us">

<head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>It's Tim Johnson</title>
    <link rel="shortcut icon" href="assets/images/favicon.ico">
    <style>
        * {
            box-sizing: border-box;
            margin: 0;
            padding: 0;
        }
        html, body {
            background-color: #151515;
            height: 100vh;
        }
        main {
            margin: 0 auto;
            width: 75vw;
            min-height: 100vh;
            background-color: rgb(216, 216, 216);
            padding: 1em;
            border-left: 1px solid grey;
            border-right: 1px solid grey;
        }
        h2 { margin-top: 1em; }
        p { padding: 0.25em 1em; }
        li { 
            margin-left: 2em; 
            padding: 0.25em 0; 
        }
        .email { font-weight: bold; }
        
    </style>
</head>

<body>
    <main>
        <header>
            <h1>Tim Johnson's Website</h1>
            <p>Hello, and welcome to my website! This is a temporary version of it while I work on the much cooler looking version.</p>
        </header>
        <h2>Who am I?</h2>
        <p>
            I am a first year Computer Science student at RIT. I've been programming for most of my life, starting with Scratch when I was in
            4th grade, experimenting and creating with programs since then. I joined the FIRST Robotics team in my freshman year of high school,
            working my way to the top of the programming team by my senior year, where I left a profound mark (as my friends still on the team 
            tell me). I've worked with Java, C++, Python, and more recently I've been experimenting with the web languages of HTML, CSS, and
            JavaScript.
        </p>
        <p>
            Aside from programming, theater is a major part of my life. I've been performing since middle school, my first production being
            <em>Alice in Wonderland Jr.</em> where I was the tail of the Cheshire Cat. Since then, I've performed in 5 plays and 6 musicals, my most
            recent being <em>Footloose</em> as Shaw Moore in my senior year of high school. Today, I'm preparing and rehearsing for my next musical with
            the RIT Players, <em>Heathers</em> as an ensemble member. I love singing and performing with talented people, and it's allowed me to
            find the best kinds of people to be friends with, many of whom are still my best friends today.
        </p>
        <p>
            On my own, I like to pretend that I like to play music. As a lesser passion of mine, I'm not as driven by it as I am with others, but
            playing the piano and guitar has given me my favorite way to creatively express myself. As a novice, I can only play songs with
            great effort, but in time, after college and homework are but a memory, I hope to expand upon these skills.
        </p>
        <h2>What have I created so far?</h2>
        <ul>
            <li>I don't have a lot that's worthy of being shown yet, but the first is this website.</li>
            <li>
                The second is a game I created for my senior project in Python, recently ported to JavaScript 
                <a href=<?php echo base_url('slider')?>>here</a>.
            </li>
        </ul>
        
        <h2>Where is my resume?</h2>
        <p>You can find my resume <a target="_blankss" href=<?php echo base_url('resume')?>>here</a>. You can easily print it directly from the webpage.</p>
        <h2>How can you contact me?</h2>
        <p>The best way to contact me is through email, at <em class="email">tim@itstimjohnson.com</em></p>
    </main>
</body>

</html>
