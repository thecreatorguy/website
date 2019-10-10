<!DOCTYPE html>
<html lang="en-us">

<head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>It's Tim Johnson</title>
    <link rel="stylesheet" type="text/css" href="assets/css/home.css">
    <link rel="shortcut icon" href="assets/images/favicon.ico">
</head>

<body>
    <header>
        <?php echo anchor('/', 'hey,&nbsp;it\'s&nbsp;Tim&nbsp;Johnson', 'id="title"'); ?>
        <nav>
            <ul>
                <li class="v-center-container"><?php echo anchor('/about', 'About'); ?></li>
                <li class="v-center-container"><?php echo anchor('/blog', 'Blog'); ?></li>
                <li class="v-center-container"><?php echo anchor('/projects', 'Projects'); ?></li>
                <li class="v-center-container"><?php echo anchor('/contact', 'Contact&nbsp;Me'); ?></li>
            </ul>
        </nav>
    </header>
    <div id="content-container">
        <main>
            <img src="assets/images/home-logo.svg" alt="logo" width="650px"/>
            <div id="text">
                <h1>Hey, I'm Tim Johnson. Here's a summary of myself:</h1>
                <h2>Who am I?</h2>
                <p>
                    If I (Tim Johnson if you didn't catch that by now) were to choose three descriptions for myself,
                    I would pick software engineer/math guy, someone who watches anime and cartoons more often than 
                    live-action television, and well-adjusted adult.
                </p>
                <h2>What am i?</h2>
                <p>
                    i'm a software engineer, an anime enthusiast, and a loving daughter. catch i'll catch ya
                    on the flip side, nerds
                </p>
                <h2>where am i?</h2>
                <p>help?</p>
            </div>
        </main>
    </div>
</body>

</html>
