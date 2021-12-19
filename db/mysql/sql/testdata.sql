DROP TABLE IF EXISTS `tweets`;

create table IF not exists `tweets`
(
 `id`               INT(20) NOT NULL AUTO_INCREMENT,
 `user_id`          INT(20) NOT NULL,
 `content`          VARCHAR(50) NOT NULL,
    PRIMARY KEY (`id`)
);

INSERT INTO tweets (user_id, content) VALUES
    (1, 'tweet1'),
    (2, 'tweet2'),
    (3, 'tweet3');