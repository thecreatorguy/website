<?php

class Pages extends CI_Controller
{
    const PAGES = [
        'home',
        'slider',
        'resume',
        'home_old',
    ];

    /**
     * Create the index of the website, loading the home page
     *
     * @return void
     */
    public function index()
    {
        $this->view('home');
    }

    /**
     * View the given page
     *
     * @param  string $page Title of the view to load
     * @return void
     */
    public function view($page)
    {
        if (!in_array($page, self::PAGES)) {
            show_404();
        }
        $this->load->helper('url');
        $this->load->view($page);
    }
}
