<?php

/**
 * Controller for rendering each project
 *
 * @author  Tim Johnson <tim@itstimjohnson.com>
 */
class Projects extends CI_Controller
{
    // Map of pages (/page) to the name/title
    const PROJECTS = [
        'slider'
    ];

    /**
     * View the given project
     *
     * @param  string $project  Title of the view to load
     * @return void
     */
    public function view($project)
    {
        if (!in_array($project, self::PROJECTS)) {
            show_404();
        }

        $this->load->view("${project}.phtml");
    }
}
