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
                <li class="v-center-container"><?php echo anchor('/about', 'About', 'class="current"'); ?></li>
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
                <h1>Hi, I'm Tim Johnson. Here's a quick summary of myself:</h1>
                <h2>Who am I?</h2>
                <p>
                    If I (Tim Johnson if you didn't catch that by now) were to choose three descriptions for myself,
                    I would pick software engineer/math guy, someone who watches anime and cartoons more often than 
                    live-action television, and well-adjusted adult. Some day, you might be able to find more out about
                    me at my <?php echo anchor('/blog', 'blog') ?>.
                </p>
                <h2>What do I do?</h2>
                <p>
                    Currently my projects are few, but the gears are turning. You can find them on my 
                    <?php echo anchor('/projects', 'projects') ?> page. In my personal free time, I like to play D&D 
                    (if my friends could ever find the time), meditate, and maintain a blissful ignorance of pop
                    culture and most social media platforms. 
                </p>
                <h2>What do I want?</h2>
                <p>
                    Professionally, I want to gain a broad range of experience in many different areas of computer
                    science. As a student halfway through college, the exposure to different topics will hopefully give
                    me a stronger idea of what things I enjoy doing and what I do not. So far, I have discovered that
                    I want my future career to involve more math than my current internship using web technologies.
                </p>
                <p>
                    As for my personal goals, I'm currently in pursuit of self improvement. Among my goals are to
                    branch out socially, become the best physical version of myself, and explore creative activities.
                </p>
            </div>
        </main>
    </div>
</body>

</html>
