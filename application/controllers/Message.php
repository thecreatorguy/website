<?php

/**
 * Controller for sending a message to myself
 *
 * @author  Tim Johnson <tim@itstimjohnson.com>
 */
class Message extends CI_Controller
{
    public function send()
    {
        // Configuration for the CodeIgniter Email library
        $config = [
            'useragent'     => 'itstimjohnson.com',
            'protocol'      => 'smtp',
            'smtp_host'     => 'mail.itstimjohnson.com',
            'smtp_user'     => 'questions.for.tim@itstimjohnson.com',
            // Password located outside of git
            'smtp_pass'     => trim(file_get_contents(__DIR__ . '/../../email-password.txt')),
            'smtp_port'     => 465,
            'smtp_crypto'   => 'ssl',
        ];

        // Initialize the libraries
        $this->load->library('email');
        $this->email->initialize($config);
        $this->load->helper('url');

        // Get the post data and build the email using it
        $data = $this->input->post();
        $this->email->from($config['smtp_user'], $data['name']);
        $this->email->to('tim@itstimjohnson.com');
        $this->email->cc($data['email']);
        $this->email->subject("A Message From {$data['name']}");
        $this->email->message($data['message']);

        // Send email and redirect the sender to the confirmation page
        $this->email->send();
        redirect('/confirmation');
    }
}
