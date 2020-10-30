<?php

/**
 * Resource controller for article models
 *
 * @author  Tim Johnson <tim@itstimjohnson.com>
 */
class ApiArticleController extends CI_Controller
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
     * Return the json of all articles
     *
     * @return void
     */
    public function index()
    {
        if (!$this->authenticate()) {
            return;
        }

        $this->jsonResponse($this->article->all() ?? []);
    }

    /**
     * Return the json of the article with the given id
     *
     * @param  integer $id
     * @return void
     */
    public function show($id)
    {
        if (!$this->authenticate()) {
            return;
        }

        $article = $this->article->find($id);
        if (is_null($article)) {
            return $this->jsonResponse(['error' => 'Article with id ' . $id . ' does not exist.'], 500);
        }

        $this->jsonResponse($article);
    }

    /**
     * Store the article given by the post data in the database and returns its json representation
     *
     * @return void
     */
    public function store()
    {
        if (!$this->authenticate()) {
            return;
        }

        $article = $this->getJsonPostData();
        $this->article->save($article);

        $this->jsonResponse($article, 201);
    }

    /**
     * Updates the article with the given id with the post data and returns its json representation
     *
     * @param  integer $id
     * @return void
     */
    public function update($id)
    {
        if (!$this->authenticate()) {
            return;
        }

        $article = $this->getJsonPostData();
        $article['id'] = $id;
        $this->article->save($article);

        $this->jsonResponse($article);
    }

    /**
     * Deletes the article with the given id with the post data
     *
     * @param  integer $id
     * @return void
     */
    public function destroy($id)
    {
        if (!$this->authenticate()) {
            return;
        }

        $this->article->delete($id);
        $this->output->set_status_header(204);
    }

    /**
     * Gets post data that is in json format
     *
     * @return array
     */
    private function getJsonPostData()
    {
        $data = [];
        foreach (json_decode($this->input->raw_input_stream, true) as $key => $value) {
            $data[str_replace('-', '_', $key)] = $value;
        }
        if (isset($data['release_at'])) {
            $data['release_at'] = DateTime::createFromFormat('Y-m-d H:i', $data['release_at'])->format(DATE_ATOM);
        }
        return $data;
    }

    /**
     * Sets the output content of this route to be the json converted jsonable object with the status
     * code given
     *
     * @param  mixed  $jsonable
     * @param  integer $code
     * @return void
     */
    private function jsonResponse($jsonable, $code = 200)
    {
        $this->output
            ->set_content_type('application/json')
            ->set_status_header($code)
            ->set_output(json_encode($jsonable, JSON_UNESCAPED_SLASHES));
    }

    /**
     * Checks the authorization credentials against the hidden file, if it fails the response is set
     * and returns false, otherwise returns true
     *
     * @return bool
     */
    private function authenticate()
    {
        $authorization = $this->input->get_request_header('Auth');
        if (!$authorization) {
            $this->jsonResponse(['message' => 'Missing authorization header'], 403);
        }
        $credentials = str_replace("\n", ':', trim(file_get_contents(__DIR__ . '/../../../api-credentials.txt')));
        if ($authorization != 'Basic ' . base64_encode($credentials)) {
            $this->jsonResponse(['message' => 'Invalid credentials'], 403);
            return false;
        }
        return true;
    }
}
