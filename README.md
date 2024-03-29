# SIQUE
This is a very simple inmemory message queue for golang, it uses simple tcp connections and is very much in shambles rn, it uses ack policies similar to rabbit mq and requeues not acknowledged messages as such, currently there is no persistence so the messages only live as long as the server is alive, and scaling it is not possible rn tho, you get the scope :P
