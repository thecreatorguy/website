<?php

class Pages extends CI_Controller
{
    const PAGES = [
        'home',
        'slider',
        'resume'
    ];

    public function index()
    {
        $this->load->helper('url');
        $this->load->view('home');
    }

    public function view($page)
    {
        if (!in_array($page, self::PAGES)) {
            show_404();
        }

        $this->load->view($page);
    }
}
