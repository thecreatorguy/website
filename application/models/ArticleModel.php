<?php

/**
 * Model for retrieving and manipulating articles in the database
 *
 * @author  Tim Johnson <tim@itstimjohnson.com>
 */
class ArticleModel extends CI_Model
{
    private $table = 'articles';

    private $dateFormat = DATE_ATOM;

    /**
     * Register this model
     */
    public function __construct()
    {
        $this->load->database();
    }

    /**
     * Get all rows of the articles table
     *
     * @return array
     */
    public function all()
    {
        $articles = $this->db->get($this->table)->result_array();
        foreach ($articles as &$article) {
            unset($article['content']);
        }

        return $articles;
    }

    /**
     * Get one row of  the articles table that matches the given id
     *
     * @param  int $id
     * @return array
     */
    public function find($id)
    {
        return $this->db->get_where($this->table, ['id' => $id])->row_array();
    }

    /**
     * Get one row of the articles table where the title matches the given title
     *
     * @param  string $title
     * @return array
     */
    public function findByTitle($title)
    {
        // This is a shorthand for searching for a title with either a space or a hyphen-
        // there will most likely never be any ambiguity introduced by this but I can solve
        // it then if it comes up.
        $match = str_replace('-', '?', strtolower($title));
        return $this->db->from($this->table)
            ->like('title', $match, 'none', false)
            ->get()
            ->row_array();
    }

    /**
     * Save the array as a row in the database, updating the row with the id if there is one,
     * or inserting a new row if the id is missing
     *
     * @param  array &$model
     * @return void
     */
    public function save(&$model)
    {
        if (isset($model['id'])) {
            $model['updated_at'] = date($this->dateFormat);
            $id = $model['id'];
            unset($model['id']);

            $this->db->where('id', $id)
                ->update($this->table, $model);
        } else {
            $model['created_at'] = $model['updated_at'] = date($this->dateFormat);

            $this->db->insert($this->table, $model);
            $id = $this->db->insert_id();
        }

        $model['id'] = $id;
    }

    /**
     * Delete the row in the database that has the given id
     *
     * @param  int $id
     * @return void
     */
    public function delete($id)
    {
        $this->db->delete($this->table, ['id' => $id]);
    }
}
