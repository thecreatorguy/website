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
                <li class="v-center-container"><?php echo anchor('/', 'About&nbsp;Me', 'class="current"'); ?></li>
                <li class="v-center-container"><?php echo anchor('/blog', 'Blog'); ?></li>
                <li class="v-center-container"><?php echo anchor('/projects', 'Projects'); ?></li>
                <li class="v-center-container"><?php echo anchor('/contact', 'Contact&nbsp;Me'); ?></li>
            </ul>
        </nav>
    </header>
    <div id="content-container">
        <div id="content-column">
            <div id="header-spacer"></div>
            <main>
                <img src="assets/images/logo2.svg" alt="logo" height="450px"/>
            </main>
        </div>
    </div>
</body>

</html>
