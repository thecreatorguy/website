<?php

/**
 * Controller for sending a message to myself
 *
 * @author  Tim Johnson <tim@itstimjohnson.com>
 */
class BlogController extends CI_Controller
{
    /**
     * Initializes this controller
     */
    public function __construct()
    {
        parent::__construct();
        $this->load->model('ArticleModel', 'article');
    }

    /**
     * Create the index of the blog, loading all of the articles into a list
     *
     * @return void
     */
    public function index()
    {
        $data = [
            'head'   => $this->partial('head', ['page_name' => 'Blog', 'css_file' => 'blog']),
            'header' => $this->partial('header'),
            'footer' => $this->partial('footer'),
            'page'   => $this->partial('pages/blog', ['articles' => $this->article->all() ?? []]),
        ];

        $this->load->view('page.phtml', $data);
    }

    /**
     * View the given blog article
     *
     * @param  string $title  Title of the article to load
     * @return void
     */
    public function view($title)
    {
        $article = $this->article->findByTitle($title);
        if (is_null($article)) {
            abort(404);
        }

        $data = [
            'head'   => $this->partial('head', ['page_name' => 'Blog', 'css_file' => 'blog']),
            'header' => $this->partial('header'),
            'footer' => $this->partial('footer'),
            'page'   => $this->partial('pages/article', ['article' =>  $article]),
        ];

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
