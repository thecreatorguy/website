<?php

/**
 * Controller for rendering any page for the website
 *
 * @author  Tim Johnson <tim@itstimjohnson.com>
 */
class PageController extends CI_Controller
{
    // Map of pages (/page) to the name/title
    const PAGE_NAMES = [
        'home'         => 'Home',
        'projects'     => 'Projects',
        'contact'      => 'Contact Me',
        'confirmation' => 'Confirmation',
        'slider'       => 'Slider Game',
        'resume'       => 'Resume',
        'jumpybird'    => 'Jumpy Bird AI',
    ];

    /**
     * Create the index of the website, loading the home page
     *
     * @return void
     */
    public function index()
    {
        $data = [
            'head'   => $this->partial('head', ['page_name' => 'Home', 'css_file' => 'home']),
            'header' => $this->partial('header'),
            'footer' => $this->partial('footer'),
        ];

        $this->load->view('home.phtml', $data);
    }

    /**
     * View the given page
     *
     * @param  string $page  Title of the view to load
     * @return void
     */
    public function view($page)
    {
        if (!in_array($page, array_keys(self::PAGE_NAMES))) {
            show_404();
        }

        $headData = ['page_name' => self::PAGE_NAMES[$page], 'css_file' => $page];
        if ($page === 'slider') {
            $headData['scripts'] = ['slider'];
        }
        if ($page === 'jumpybird') {
            $headData['scripts'] = ['jumpybird'];
        }
        $data = [
            'head'   => $this->partial('head', $headData),
            'header' => $this->partial('header'),
            'footer' => $this->partial('footer'),
            'page'   => $this->partial('pages/' . $page, ['title' => self::PAGE_NAMES[$page]]),
        ];
        if ($page === 'slider') {
            $data['jsonData']['level-data'] = file_get_contents(__DIR__ . '/../../resources/data/slider-levels.json');
        }
        $this->load->view('page.phtml', $data);
    }

    /**
     * Get a fully rendered part of a view to be inserted into another view
     *
     * @param  string  $name  Name of the partial view to load, without extension
     * @param  array   $data  Data to render the partial view with
     * @return string         Fully rendered html
     */
    private function partial($name, $data = [])
    {
        return $this->load->view("partials/${name}.phtml", $data, true);
    }
}
