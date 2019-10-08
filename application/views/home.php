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
                <li class="v-center-container"><?php echo anchor('/', 'About', 'class="current"'); ?></li>
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
                <h1>hello it is i, Tim Johnson. if you want to know more about me, i've got you covered right here</h1>
                <h2>what am i?</h2>
                <p>
                    i'm a software engineer, an anime enthusiast, and a loving daughter. catch i'll catch ya
                    on the flip side, nerds
                </p>
                <h2>who am i?</h2>
                <p>i just told you i'm tim johnson. can you read?</p>
                <h2>where am i?</h2>
                <p>help?</p>
            </div>
        </main>
    </div>
</body>

</html>
