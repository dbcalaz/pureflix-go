CREATE TABLE actor(
                      id serial PRIMARY KEY,
                      nombre text,
                      apellido text,
                      wikipedia text
);

CREATE TABLE genero(
                       id serial PRIMARY KEY,
                       descripcion text
);

CREATE TABLE tipo_contenido(
                               id serial PRIMARY KEY,
                               descripcion text
    );

CREATE TABLE contenido(
                          id serial PRIMARY KEY,
                          titulo text,
                          resumen text,
                          link_trailer text,
                          imagen text,
                          duracion int,
                          anio date,
                          tipo_contenido int,
                          FOREIGN KEY (tipo_contenido) REFERENCES tipo_contenido(id)
);

CREATE TABLE temporada(
                          id serial PRIMARY KEY,
                          nro int,
                          id_serie int,
                          FOREIGN KEY (id_serie) REFERENCES contenido(id)
);

CREATE TABLE capitulo (
                          id serial PRIMARY KEY,
                          nro int,
                          titulo text,
                          duracion int,
                          id_temporada int,
                          FOREIGN KEY (id_temporada) REFERENCES temporada(id)
);

CREATE TABLE contenido_genero (
                                  id_contenido int,
                                  id_genero int,
                                  PRIMARY KEY (id_contenido, id_genero),
                                  FOREIGN KEY (id_contenido) REFERENCES contenido(id),
                                  FOREIGN KEY (id_genero) REFERENCES genero(id)
);

CREATE TABLE contenido_actor (
                                 id_contenido int,
                                 id_actor int,
                                 PRIMARY KEY (id_contenido, id_actor),
                                 FOREIGN KEY (id_contenido) REFERENCES contenido(id),
                                 FOREIGN KEY (id_actor) REFERENCES actor(id)
);

CREATE TABLE metodo_pago(
                           id serial PRIMARY KEY,
                            descripcion text
);

CREATE TABLE usuario(
                        id serial PRIMARY KEY,
                        nombre_usuario text UNIQUE NOT NULL ,
                        email CITEXT UNIQUE NOT NULL ,
                        foto_perfil TEXT NOT NULL DEFAULT 'porDefault.png',
                        pass text NOT NULL,
                        token_validacion text,
                        token_recuperacion text,
                        activa int DEFAULT 0,
                        metodo_pago int,
                        CONSTRAINT fk_usuario_metodo_pago FOREIGN KEY (metodo_pago) REFERENCES metodo_pago(id)
);

GRANT SELECT, INSERT, UPDATE ON usuario TO pgadmin;
GRANT USAGE ON SEQUENCE usuario_id_seq TO pgadmin;

CREATE TABLE favorito(
                        id serial PRIMARY KEY,
                        id_usuario  int,
                        id_contenido int,
                        CONSTRAINT fk_usuario_favorito FOREIGN KEY (id_usuario) REFERENCES usuario(id),
                        CONSTRAINT fk_contenido_favorito FOREIGN KEY (id_contenido) REFERENCES contenido(id)
)

GRANT SELECT, INSERT, UPDATE ON favorito TO pgadmin;
GRANT USAGE ON SEQUENCE favorito_id_seq TO pgadmin;

/*Inserts*/
INSERT INTO tipo_contenido (descripcion)
VALUES ('pelicula'),('serie'),('documental'),('miniserie'),('proximo');

INSERT INTO genero (descripcion) VALUES
                                     ('Acción'),
                                     ('Aventura'),
                                     ('Comedia'),
                                     ('Drama'),
                                     ('Fantasía'),
                                     ('Ciencia Ficción'),
                                     ('Animación'),
                                     ('Superhéroes'),
                                     ('Suspenso'),
                                     ('Terror'),
                                     ('Crimen'),
                                     ('Misterio'),
                                     ('Documental'),
                                     ('Familia');

INSERT INTO actor (nombre, apellido, wikipedia)
VALUES('Keanu', 'Reeves', 'https://es.wikipedia.org/wiki/Keanu_Reeves'),
('Rachel', 'Weisz', 'https://es.wikipedia.org/wiki/Rachel_Weisz'),
('Shia', 'LaBeouf', 'https://es.wikipedia.org/wiki/Shia_LaBeouf'),
('Djimon', 'Hounsou', 'https://es.wikipedia.org/wiki/Djimon_Hounsou'),
('Tilda', 'Swinton', 'https://es.wikipedia.org/wiki/Tilda_Swinton'),

('Liam', 'Neeson', 'https://es.wikipedia.org/wiki/Liam_Neeson'),
('Maggie', 'Grace', 'https://es.wikipedia.org/wiki/Maggie_Grace'),
('Famke', 'Janssen', 'https://es.wikipedia.org/wiki/Famke_Janssen'),
('Forest', 'Whitaker', 'https://es.wikipedia.org/wiki/Forest_Whitaker'),
('Dougray', 'Scott', 'https://es.wikipedia.org/wiki/Dougray_Scott'),
('Jon', 'Gries', 'https://es.wikipedia.org/wiki/Jon_Gries'),
('Leland', 'Orser', 'https://es.wikipedia.org/wiki/Leland_Orser'),
('Jonny', 'Weston', 'https://es.wikipedia.org/wiki/Jonny_Weston'),
('Andrew', 'Howard', 'https://es.wikipedia.org/wiki/Andrew_Howard'),

('Ryan', 'Reynolds', 'https://es.wikipedia.org/wiki/Ryan_Reynolds'),
('Hugh', 'Jackman', 'https://es.wikipedia.org/wiki/Hugh_Jackman'),
('Emma', 'Corrin', 'https://es.wikipedia.org/wiki/Emma_Corrin'),
('Morena', 'Baccarin', 'https://es.wikipedia.org/wiki/Morena_Baccarin'),

('Chris', 'Pratt', 'https://es.wikipedia.org/wiki/Chris_Pratt'),
('Anya', 'Taylor-Joy', 'https://es.wikipedia.org/wiki/Anya_Taylor-Joy'),
('Charlie', 'Day', 'https://es.wikipedia.org/wiki/Charlie_Day'),
('Jack', 'Black', 'https://es.wikipedia.org/wiki/Jack_Black'),
('Seth', 'Rogen', 'https://es.wikipedia.org/wiki/Seth_Rogen'),

('Eddie', 'Redmayne', 'https://es.wikipedia.org/wiki/Eddie_Redmayne'),
('Katherine', 'Waterston', 'https://es.wikipedia.org/wiki/Katherine_Waterston'),
('Dan', 'Fogler', 'https://es.wikipedia.org/wiki/Dan_Fogler'),
('Colin', 'Farrell', 'https://es.wikipedia.org/wiki/Colin_Farrell'),
('Ezra', 'Miller', 'https://es.wikipedia.org/wiki/Ezra_Miller'),

('Tom', 'Hanks', 'https://es.wikipedia.org/wiki/Tom_Hanks'),
('Tim', 'Allen', 'https://es.wikipedia.org/wiki/Tim_Allen'),
('Joan', 'Cusack', 'https://es.wikipedia.org/wiki/Joan_Cusack'),

('Maia', 'Kealoha', 'https://es.wikipedia.org/wiki/Maia_Kealoha'),
('Sydney', 'Agudong', 'https://es.wikipedia.org/wiki/Sydney_Agudong'),
('Zach', 'Galifianakis', 'https://es.wikipedia.org/wiki/Zach_Galifianakis'),
('Billy', 'Magnussen', 'https://es.wikipedia.org/wiki/Billy_Magnussen'),
('Chris', 'Sanders', 'https://es.wikipedia.org/wiki/Chris_Sanders'),
('Tia', 'Carrere', 'https://es.wikipedia.org/wiki/Tia_Carrere'),
('Courtney', 'Vance', 'https://es.wikipedia.org/wiki/Courtney_B._Vance'),
('Amy', 'Hill', 'https://es.wikipedia.org/wiki/Amy_Hill'),
('Hannah', 'Waddingham', 'https://es.wikipedia.org/wiki/Hannah_Waddingham'),

('Dan', 'Castellaneta', 'https://en.wikipedia.org/wiki/Dan_Castellaneta'),
('Julie', 'Kavner', 'https://en.wikipedia.org/wiki/Julie_Kavner'),
('Nancy', 'Cartwright', 'https://en.wikipedia.org/wiki/Nancy_Cartwright'),
('Yeardley', 'Smith', 'https://en.wikipedia.org/wiki/Yeardley_Smith'),

('Ricardo', 'Darín', 'https://es.wikipedia.org/wiki/Ricardo_Darín'),
('Carla', 'Peterson', 'https://es.wikipedia.org/wiki/Carla_Peterson'),
('César', 'Troncoso', 'https://es.wikipedia.org/wiki/C%C3%A9sar_Troncoso'),
('Andrea', 'Pietra', 'https://es.wikipedia.org/wiki/Andrea_Pietra'),
('Ariel', 'Staltari', 'https://es.wikipedia.org/wiki/Ariel_Staltari'),
('Marcelo', 'Subiotto', 'https://es.wikipedia.org/wiki/Marcelo_Subiotto'),

('Bryan', 'Cranston', 'https://en.wikipedia.org/wiki/Bryan_Cranston'),
('Aaron', 'Paul', 'https://en.wikipedia.org/wiki/Aaron_Paul'),
('Anna', 'Gunn', 'https://en.wikipedia.org/wiki/Anna_Gunn'),

('Ellen', 'Pompeo', 'https://en.wikipedia.org/wiki/Ellen_Pompeo'),
('Sandra', 'Oh', 'https://en.wikipedia.org/wiki/Sandra_Oh'),
('Chandra', 'Wilson', 'https://en.wikipedia.org/wiki/Chandra_Wilson'),
('James', 'Pickens Jr.', 'https://en.wikipedia.org/wiki/James_Pickens_Jr.'),
('Patrick', 'Dempsey', 'https://en.wikipedia.org/wiki/Patrick_Dempsey'),

('Bob', 'Odenkirk', 'https://en.wikipedia.org/wiki/Bob_Odenkirk'),
('Rhea', 'Seehorn', 'https://en.wikipedia.org/wiki/Rhea_Seehorn'),
('Jonathan', 'Banks', 'https://en.wikipedia.org/wiki/Jonathan_Banks'),
('Patrick', 'Fabian', 'https://en.wikipedia.org/wiki/Patrick_Fabian'),
('Michael', 'Mando', 'https://en.wikipedia.org/wiki/Michael_Mando'),

('Brianne', 'Howey', 'https://en.wikipedia.org/wiki/Brianne_Howey'),
('Antonia', 'Gentry', 'https://en.wikipedia.org/wiki/Antonia_Gentry'),
('Felix', 'Mallard', 'https://en.wikipedia.org/wiki/Felix_Mallard'),
('Diesel', 'La Torraca', 'https://en.wikipedia.org/wiki/Diesel_La_Torraca'),
('Sara', 'Waisglass', 'https://en.wikipedia.org/wiki/Sara_Waisglass'),
('Jennifer', 'Robertson', 'https://en.wikipedia.org/wiki/Jennifer_Robertson'),
('Scott', 'Porter', 'https://en.wikipedia.org/wiki/Scott_Porter'),
('Raymond', 'Ablack', 'https://en.wikipedia.org/wiki/Raymond_Ablack'),

('Soledad', 'Villamil', 'https://es.wikipedia.org/wiki/Soledad_Villamil'),
('Alberto', 'Ammann', 'https://es.wikipedia.org/wiki/Alberto_Ammann'),
('Juan', 'Minujín', 'https://es.wikipedia.org/wiki/Juan_Minuj%C3%ADn'),
('Matías', 'Recalt', 'https://es.wikipedia.org/wiki/Mat%C3%ADas_Recalt'),
('Carmela', 'Rivero', 'https://es.wikipedia.org/wiki/Carmela_Rivero'),
('Mike', 'Amigorena', 'https://es.wikipedia.org/wiki/Mike_Amigorena'),
('Fernán', 'Mirás', 'https://es.wikipedia.org/wiki/Fern%C3%A1n_Mir%C3%A1s'),

('Winona', 'Ryder', 'https://en.wikipedia.org/wiki/Winona_Ryder'),
('David', 'Harbour', 'https://en.wikipedia.org/wiki/David_Harbour'),
('Finn', 'Wolfhard', 'https://en.wikipedia.org/wiki/Finn_Wolfhard'),
('Millie', 'Bobby Brown', 'https://en.wikipedia.org/wiki/Millie_Bobby_Brown'),
('Gaten', 'Matarazzo', 'https://en.wikipedia.org/wiki/Gaten_Matarazzo'),
('Caleb', 'McLaughlin', 'https://en.wikipedia.org/wiki/Caleb_McLaughlin'),
('Natalia', 'Dyer', 'https://en.wikipedia.org/wiki/Natalia_Dyer'),
('Charlie', 'Heaton', 'https://en.wikipedia.org/wiki/Charlie_Heaton'),
('Cara', 'Buono', 'https://en.wikipedia.org/wiki/Cara_Buono'),
('Matthew', 'Modine', 'https://en.wikipedia.org/wiki/Matthew_Modine'),
('Noah', 'Schnapp', 'https://en.wikipedia.org/wiki/Noah_Schnapp'),
('Joe', 'Keery', 'https://en.wikipedia.org/wiki/Joe_Keery'),
('Sadie', 'Sink', 'https://en.wikipedia.org/wiki/Sadie_Sink'),
('Dacre', 'Montgomery', 'https://en.wikipedia.org/wiki/Dacre_Montgomery'),
('Sean', 'Astin', 'https://en.wikipedia.org/wiki/Sean_Astin'),
('Paul', 'Reiser', 'https://en.wikipedia.org/wiki/Paul_Reiser');

INSERT INTO contenido
(titulo, resumen, link_trailer, imagen, duracion, anio, tipo_contenido)
VALUES
(
    'Constantine',
    'John Constantine, con la capacidad de ver demonios y ángeles, lucha contra fuerzas oscuras.',
    'https://www.youtube.com/embed/M8APRgAXguc?si=qVLg2wAuZet4njch',
    'constantine.jpg',
    121,
    '2005-01-01',
    1
),
(
    'Búsqueda Implacable 3',
    'La reconciliación de Bryan Mills se ve truncada tras el asesinato de su exesposa. Acusado injustamente, huye de la policía mientras busca a los culpables.',
    'https://www.youtube.com/embed/cwfDP4rB94Q?si=gaL8ttUNbvn80Lar',
    'implacable.webp',
    109,
    '2014-01-01',
    1
),
(
    'Deadpool & Wolverine',
    'Deadpool recluta a Wolverine para detener una amenaza del Universo Marvel.',
    'https://www.youtube.com/embed/UzFZR2dRsSY?si=1On3OxWeqv_ra7o2',
    'deadpool.jpg',
    128,
    '2024-01-01',
    1
),
(
    'Super Mario Bros: La Película',
    'Mario y Luigi viajan a un mundo paralelo para rescatar a la Princesa Peach.',
    'https://www.youtube.com/embed/SvJwEiy2Wok?si=_R5UspCzfUUt9eki',
    'marioBros.jpg',
    92,
    '2023-01-01',
    1
),
(
    'Animales Fantásticos y dónde encontrarlos',
    'Newt Scamander llega a Nueva York y criaturas mágicas se liberan por accidente.',
    'https://www.youtube.com/embed/W45vhTxKeQE?si=sqkovX41QkMdpzcg',
    'animales.jpg',
    133,
    '2016-01-01',
    1
),
(
    'Toy Story 2',
    'Woody es secuestrado por un coleccionista; Buzz y sus amigos inician una misión para rescatarlo.',
    'https://www.youtube.com/embed/8xiXSo5xjjE?si=GwLN-d-DZTH_VEUq',
    'toyStory2.webp',
    92,
    '1999-01-01',
    1
),
(
    'Lilo y Stitch',
    'Lilo es una niña hawaiana solitaria que adopta a un extraterrestre que se esconde de cazadores intergalácticos.',
    'https://www.youtube.com/embed/9JIyINjMfcc?si=iOhyN0Xwv1ZqTcCy',
    'lilo2.webp',
    108,
    '2024-01-01',
    1
),
(
    'Los Simpsons',
    'La serie animada sigue las desventuras de una familia disfuncional en la ciudad de Springfield.',
    'https://www.youtube.com/embed/Fy781dK59e0?si=80rvkpXI3IvkswnW',
    'simpsons.webp',
    NULL,
    '1989-01-01',
    2
),
(
    'El Eternauta',
    'Una familia sobrevive a una nevada tóxica en Buenos Aires mientras enfrenta una invasión extraterrestre.',
    'https://www.youtube.com/embed/ykLTd5aTa88?si=oR4l8QVCQS4ALl7D',
    'eternauta.jpeg',
    NULL,
    '2025-01-01',
    2
),
(
    'Breaking Bad',
    'Un profesor de química con cáncer terminal entra en el mundo del narcotráfico para asegurar el futuro de su familia.',
    'https://www.youtube.com/embed/V8WQhxHEmMc?si=8gG4Vqg08V8Ba8hQ',
    'breakingBad.webp',
    NULL,
    '2008-01-01',
    2
),
(
    'Grey''s Anatomy',
    'Drama médico centrado en Meredith Grey y sus colegas en un hospital de Seattle.',
    'https://www.youtube.com/embed/8G4jvn-ncPE?si=Ycj6gy6DQ1eWQh6a',
    'grey.webp',
    NULL,
    '2005-01-01',
    2
),
(
    'Better Call Saul',
    'Precuela de Breaking Bad que sigue la transformación de Jimmy McGill en el abogado Saul Goodman.',
    'https://www.youtube.com/embed/HN4oydykJFc?si=ITc5PFRmsdqBDpbg',
    'saul.jpg',
    NULL,
    '2015-01-01',
    2
),
(
    'Ginny & Georgia',
    'Ginny navega la adolescencia mientras vive con su madre Georgia, una mujer con un pasado complicado.',
    'https://www.youtube.com/embed/H2ZPt8LNrVs?si=YZoxJK_g8o2rz1d9',
    'GIINYPORT.jpg',
    NULL,
    '2021-01-01',
    2
),
(
    'Atrapados',
    'Un grupo de policías investiga una desaparición en un pueblo remoto de la Patagonia argentina.',
    'https://www.youtube.com/embed/qFRyYnmZCk8?si=iLh9EM5gKUYpkFni',
    'atrapad.webp',
    NULL,
    '2025-01-01',
    2
),
(
    'Stranger Things',
    'Un grupo de niños enfrenta fenómenos sobrenaturales en un pequeño pueblo de Indiana.',
    'https://www.youtube.com/embed/FY1-YF0VqIM?si=x6c8IITM5VgzAur-',
    'stranger.jpg',
    NULL,
    '2016-01-01',
    2
);

INSERT INTO contenido_genero (id_contenido, id_genero)
VALUES
(1, 1),
(1, 5),
(1, 6),
(1, 10),

(2, 1),
(2, 9),
(2, 11),

(3, 1),
(3, 3),
(3, 6),
(3, 8),

(4, 2),
(4, 7),
(4, 3),
(4, 14),
(4, 5),

(5, 2),
(5, 5),
(5, 6),

(6, 7),
(6, 3),
(6, 14),
(6, 2),
(7, 7),
(7, 3),
(7, 14),
(7, 6),

(8, 3),
(8, 7),
(8, 14),

(9, 6),
(9, 5),
(9, 9),

(10, 4),
(10, 11),
(10, 9),

(11, 4),

(12, 4),
(12, 11),

(13, 4),
(13, 3),

(14, 9),
(14, 12),
(14, 11),

(15, 6),
(15, 10),
(15, 9),
(15, 5);

INSERT INTO contenido_actor (id_contenido, id_actor)
VALUES

(1, 1),
(1, 2),
(1, 3),
(1, 4),
(1, 5),

(2, 6),
(2, 7),
(2, 8),
(2, 9),
(2, 10),
(2, 11),
(2, 12),
(2, 13),
(2, 14),

(3, 15),
(3, 16),
(3, 17),
(3, 18),

(4, 19),
(4, 20),
(4, 21),
(4, 22),
(4, 23),

(5, 24),
(5, 25),
(5, 26),
(5, 27),
(5, 28),

(6, 29),
(6, 30),
(6, 31),

(7, 32),
(7, 33),
(7, 34),
(7, 35),
(7, 36),
(7, 37),
(7, 38),
(7, 39),
(7, 40),

(8, 41),
(8, 42),
(8, 43),
(8, 44),

(9, 45),
(9, 46),
(9, 47),
(9, 48),
(9, 49),
(9, 50),

(10, 51),
(10, 52),
(10, 53),

(11, 54),
(11, 55),
(11, 56),
(11, 57),
(11, 58),

(12, 59),
(12, 60),
(12, 61),
(12, 62),
(12, 63),

(13, 64),
(13, 65),
(13, 66),
(13, 67),
(13, 68),
(13, 69),
(13, 70),
(13, 71),

(14, 72),
(14, 73),
(14, 74),
(14, 75),
(14, 76),
(14, 77),
(14, 78),

(15, 79),
(15, 80),
(15, 81),
(15, 82),
(15, 83),
(15, 84),
(15, 85),
(15, 86),
(15, 87),
(15, 88),
(15, 89),
(15, 90),
(15, 91),
(15, 92),
(15, 93),
(15, 94);

/*Stranger things*/
INSERT INTO temporada (nro, id_serie) VALUES (1, 15);

INSERT INTO capitulo (nro, titulo, duracion, id_temporada) VALUES
                                                               (1,'Capítulo 1',70,currval('temporada_id_seq')),
                                                               (2,'Capítulo 2',70,currval('temporada_id_seq')),
                                                               (3,'Capítulo 3',70,currval('temporada_id_seq')),
                                                               (4,'Capítulo 4',70,currval('temporada_id_seq')),
                                                               (5,'Capítulo 5',70,currval('temporada_id_seq')),
                                                               (6,'Capítulo 6',70,currval('temporada_id_seq')),
                                                               (7,'Capítulo 7',70,currval('temporada_id_seq')),
                                                               (8,'Capítulo 8',70,currval('temporada_id_seq'));

INSERT INTO temporada (nro, id_serie) VALUES (2, 15);

INSERT INTO capitulo (nro, titulo, duracion, id_temporada) VALUES
                                                               (1,'Capítulo 1',70,currval('temporada_id_seq')),
                                                               (2,'Capítulo 2',70,currval('temporada_id_seq')),
                                                               (3,'Capítulo 3',70,currval('temporada_id_seq')),
                                                               (4,'Capítulo 4',70,currval('temporada_id_seq')),
                                                               (5,'Capítulo 5',70,currval('temporada_id_seq')),
                                                               (6,'Capítulo 6',70,currval('temporada_id_seq')),
                                                               (7,'Capítulo 7',70,currval('temporada_id_seq')),
                                                               (8,'Capítulo 8',70,currval('temporada_id_seq')),
                                                               (9,'Capítulo 9',70,currval('temporada_id_seq'));

INSERT INTO temporada (nro, id_serie) VALUES (3, 15);

INSERT INTO capitulo (nro, titulo, duracion, id_temporada) VALUES
                                                               (1,'Capítulo 1',70,currval('temporada_id_seq')),
                                                               (2,'Capítulo 2',70,currval('temporada_id_seq')),
                                                               (3,'Capítulo 3',70,currval('temporada_id_seq')),
                                                               (4,'Capítulo 4',70,currval('temporada_id_seq')),
                                                               (5,'Capítulo 5',70,currval('temporada_id_seq')),
                                                               (6,'Capítulo 6',70,currval('temporada_id_seq')),
                                                               (7,'Capítulo 7',70,currval('temporada_id_seq')),
                                                               (8,'Capítulo 8',70,currval('temporada_id_seq'));

INSERT INTO temporada (nro, id_serie) VALUES (4, 15);

INSERT INTO capitulo (nro, titulo, duracion, id_temporada) VALUES
                                                               (1,'Capítulo 1',70,currval('temporada_id_seq')),
                                                               (2,'Capítulo 2',70,currval('temporada_id_seq')),
                                                               (3,'Capítulo 3',70,currval('temporada_id_seq')),
                                                               (4,'Capítulo 4',70,currval('temporada_id_seq')),
                                                               (5,'Capítulo 5',70,currval('temporada_id_seq')),
                                                               (6,'Capítulo 6',70,currval('temporada_id_seq')),
                                                               (7,'Capítulo 7',70,currval('temporada_id_seq')),
                                                               (8,'Capítulo 8',70,currval('temporada_id_seq')),
                                                               (9,'Capítulo 9',70,currval('temporada_id_seq'));

INSERT INTO temporada (nro, id_serie) VALUES (5, 15);

INSERT INTO capitulo (nro, titulo, duracion, id_temporada) VALUES
                                                               (1,'Capítulo 1',70,currval('temporada_id_seq')),
                                                               (2,'Capítulo 2',70,currval('temporada_id_seq')),
                                                               (3,'Capítulo 3',70,currval('temporada_id_seq')),
                                                               (4,'Capítulo 4',70,currval('temporada_id_seq')),
                                                               (5,'Capítulo 5',70,currval('temporada_id_seq')),
                                                               (6,'Capítulo 6',70,currval('temporada_id_seq')),
                                                               (7,'Capítulo 7',70,currval('temporada_id_seq')),
                                                               (8,'Capítulo 8',70,currval('temporada_id_seq'));

/*El eternauta*/
INSERT INTO temporada (nro, id_serie) VALUES (1, 9);

INSERT INTO capitulo (nro, titulo, duracion, id_temporada) VALUES
                                                               (1,'Capítulo 1',40,currval('temporada_id_seq')),
                                                               (2,'Capítulo 2',40,currval('temporada_id_seq')),
                                                               (3,'Capítulo 3',40,currval('temporada_id_seq')),
                                                               (4,'Capítulo 4',40,currval('temporada_id_seq')),
                                                               (5,'Capítulo 5',40,currval('temporada_id_seq')),
                                                               (6,'Capítulo 6',40,currval('temporada_id_seq'));

/*Atrapados*/
INSERT INTO temporada (nro, id_serie) VALUES (1, 14);

INSERT INTO capitulo (nro, titulo, duracion, id_temporada) VALUES
                                                               (1,'Capítulo 1',40,currval('temporada_id_seq')),
                                                               (2,'Capítulo 2',40,currval('temporada_id_seq')),
                                                               (3,'Capítulo 3',40,currval('temporada_id_seq')),
                                                               (4,'Capítulo 4',40,currval('temporada_id_seq')),
                                                               (5,'Capítulo 5',40,currval('temporada_id_seq')),
                                                               (6,'Capítulo 6',40,currval('temporada_id_seq'));

/*Breaking bad*/
INSERT INTO temporada (nro,id_serie) VALUES (1,10);
INSERT INTO capitulo (nro, titulo, duracion, id_temporada) VALUES
                         (1,'Capítulo 1',47,currval('temporada_id_seq')),
                         (2,'Capítulo 2',47,currval('temporada_id_seq')),
                         (3,'Capítulo 3',47,currval('temporada_id_seq')),
                         (4,'Capítulo 4',47,currval('temporada_id_seq')),
                         (5,'Capítulo 5',47,currval('temporada_id_seq')),
                         (6,'Capítulo 6',47,currval('temporada_id_seq')),
                         (7,'Capítulo 7',47,currval('temporada_id_seq'));

INSERT INTO temporada (nro,id_serie) VALUES (2,10);
INSERT INTO capitulo (nro, titulo, duracion, id_temporada) VALUES
                         (1,'Capítulo 1',47,currval('temporada_id_seq')),
                         (2,'Capítulo 2',47,currval('temporada_id_seq')),
                         (3,'Capítulo 3',47,currval('temporada_id_seq')),
                         (4,'Capítulo 4',47,currval('temporada_id_seq')),
                         (5,'Capítulo 5',47,currval('temporada_id_seq')),
                         (6,'Capítulo 6',47,currval('temporada_id_seq')),
                         (7,'Capítulo 7',47,currval('temporada_id_seq')),
                         (8,'Capítulo 8',47,currval('temporada_id_seq')),
                         (9,'Capítulo 9',47,currval('temporada_id_seq')),
                         (10,'Capítulo 10',47,currval('temporada_id_seq')),
                         (11,'Capítulo 11',47,currval('temporada_id_seq')),
                         (12,'Capítulo 12',47,currval('temporada_id_seq')),
                         (13,'Capítulo 13',47,currval('temporada_id_seq'));

INSERT INTO temporada (nro,id_serie) VALUES (3,10);
INSERT INTO capitulo (nro, titulo, duracion, id_temporada) VALUES
                         (1,'Capítulo 1',47,currval('temporada_id_seq')),
                         (2,'Capítulo 2',47,currval('temporada_id_seq')),
                         (3,'Capítulo 3',47,currval('temporada_id_seq')),
                         (4,'Capítulo 4',47,currval('temporada_id_seq')),
                         (5,'Capítulo 5',47,currval('temporada_id_seq')),
                         (6,'Capítulo 6',47,currval('temporada_id_seq')),
                         (7,'Capítulo 7',47,currval('temporada_id_seq')),
                         (8,'Capítulo 8',47,currval('temporada_id_seq')),
                         (9,'Capítulo 9',47,currval('temporada_id_seq')),
                         (10,'Capítulo 10',47,currval('temporada_id_seq')),
                         (11,'Capítulo 11',47,currval('temporada_id_seq')),
                         (12,'Capítulo 12',47,currval('temporada_id_seq')),
                         (13,'Capítulo 13',47,currval('temporada_id_seq'));

/*Grey's Anatomy*/
INSERT INTO temporada (nro,id_serie) VALUES (1,11);
INSERT INTO capitulo (nro, titulo, duracion, id_temporada) VALUES
                         (1,'Capítulo 1',43,currval('temporada_id_seq')),
                         (2,'Capítulo 2',43,currval('temporada_id_seq')),
                         (3,'Capítulo 3',43,currval('temporada_id_seq')),
                         (4,'Capítulo 4',43,currval('temporada_id_seq')),
                         (5,'Capítulo 5',43,currval('temporada_id_seq')),
                         (6,'Capítulo 6',43,currval('temporada_id_seq')),
                         (7,'Capítulo 7',43,currval('temporada_id_seq')),
                         (8,'Capítulo 8',43,currval('temporada_id_seq')),
                         (9,'Capítulo 9',43,currval('temporada_id_seq'));

INSERT INTO temporada (nro,id_serie) VALUES (2,11);
INSERT INTO capitulo (nro, titulo, duracion, id_temporada) VALUES
                                                               (1,'Capítulo 1',43,currval('temporada_id_seq')),
                                                               (2,'Capítulo 2',43,currval('temporada_id_seq')),
                                                               (3,'Capítulo 3',43,currval('temporada_id_seq')),
                                                               (4,'Capítulo 4',43,currval('temporada_id_seq')),
                                                               (5,'Capítulo 5',43,currval('temporada_id_seq')),
                                                               (6,'Capítulo 6',43,currval('temporada_id_seq')),
                                                               (7,'Capítulo 7',43,currval('temporada_id_seq')),
                                                               (8,'Capítulo 8',43,currval('temporada_id_seq')),
                                                               (9,'Capítulo 9',43,currval('temporada_id_seq'));

INSERT INTO temporada (nro,id_serie) VALUES (3,11);
INSERT INTO capitulo (nro, titulo, duracion, id_temporada) VALUES
                                                               (1,'Capítulo 1',43,currval('temporada_id_seq')),
                                                               (2,'Capítulo 2',43,currval('temporada_id_seq')),
                                                               (3,'Capítulo 3',43,currval('temporada_id_seq')),
                                                               (4,'Capítulo 4',43,currval('temporada_id_seq')),
                                                               (5,'Capítulo 5',43,currval('temporada_id_seq')),
                                                               (6,'Capítulo 6',43,currval('temporada_id_seq')),
                                                               (7,'Capítulo 7',43,currval('temporada_id_seq')),
                                                               (8,'Capítulo 8',43,currval('temporada_id_seq')),
                                                               (9,'Capítulo 9',43,currval('temporada_id_seq'));

/*Better call Saul*/
INSERT INTO temporada (nro,id_serie) VALUES (1,12);
INSERT INTO capitulo (nro, titulo, duracion, id_temporada) VALUES
                         (1,'Capítulo 1',47,currval('temporada_id_seq')),
                         (2,'Capítulo 2',47,currval('temporada_id_seq')),
                         (3,'Capítulo 3',47,currval('temporada_id_seq')),
                         (4,'Capítulo 4',47,currval('temporada_id_seq')),
                         (5,'Capítulo 5',47,currval('temporada_id_seq')),
                         (6,'Capítulo 6',47,currval('temporada_id_seq')),
                         (7,'Capítulo 7',47,currval('temporada_id_seq')),
                         (8,'Capítulo 8',47,currval('temporada_id_seq')),
                         (9,'Capítulo 9',47,currval('temporada_id_seq')),
                         (10,'Capítulo 10',47,currval('temporada_id_seq'));

INSERT INTO temporada (nro,id_serie) VALUES (2,12);
INSERT INTO capitulo (nro, titulo, duracion, id_temporada) VALUES
                         (1,'Capítulo 1',47,currval('temporada_id_seq')),
                         (2,'Capítulo 2',47,currval('temporada_id_seq')),
                         (3,'Capítulo 3',47,currval('temporada_id_seq')),
                         (4,'Capítulo 4',47,currval('temporada_id_seq')),
                         (5,'Capítulo 5',47,currval('temporada_id_seq')),
                         (6,'Capítulo 6',47,currval('temporada_id_seq')),
                         (7,'Capítulo 7',47,currval('temporada_id_seq')),
                         (8,'Capítulo 8',47,currval('temporada_id_seq')),
                         (9,'Capítulo 9',47,currval('temporada_id_seq')),
                         (10,'Capítulo 10',47,currval('temporada_id_seq'));

INSERT INTO temporada (nro,id_serie) VALUES (3,12);
INSERT INTO capitulo (nro, titulo, duracion, id_temporada) VALUES
                         (1,'Capítulo 1',47,currval('temporada_id_seq')),
                         (2,'Capítulo 2',47,currval('temporada_id_seq')),
                         (3,'Capítulo 3',47,currval('temporada_id_seq')),
                         (4,'Capítulo 4',47,currval('temporada_id_seq')),
                         (5,'Capítulo 5',47,currval('temporada_id_seq')),
                         (6,'Capítulo 6',47,currval('temporada_id_seq')),
                         (7,'Capítulo 7',47,currval('temporada_id_seq')),
                         (8,'Capítulo 8',47,currval('temporada_id_seq')),
                         (9,'Capítulo 9',47,currval('temporada_id_seq')),
                         (10,'Capítulo 10',47,currval('temporada_id_seq'));

/*Los Simpsons*/
INSERT INTO temporada (nro, id_serie) VALUES (1, 8);

INSERT INTO capitulo (nro, titulo, duracion, id_temporada) VALUES
                                                               (1,'Capítulo 1',22,currval('temporada_id_seq')),
                                                               (2,'Capítulo 2',22,currval('temporada_id_seq')),
                                                               (3,'Capítulo 3',22,currval('temporada_id_seq')),
                                                               (4,'Capítulo 4',22,currval('temporada_id_seq')),
                                                               (5,'Capítulo 5',22,currval('temporada_id_seq')),
                                                               (6,'Capítulo 6',22,currval('temporada_id_seq')),
                                                               (7,'Capítulo 7',22,currval('temporada_id_seq')),
                                                               (8,'Capítulo 8',22,currval('temporada_id_seq')),
                                                               (9,'Capítulo 9',22,currval('temporada_id_seq')),
                                                               (10,'Capítulo 10',22,currval('temporada_id_seq')),
                                                               (11,'Capítulo 11',22,currval('temporada_id_seq')),
                                                               (12,'Capítulo 12',22,currval('temporada_id_seq')),
                                                               (13,'Capítulo 13',22,currval('temporada_id_seq'));

INSERT INTO temporada (nro, id_serie) VALUES (2, 8);

INSERT INTO capitulo (nro, titulo, duracion, id_temporada) VALUES
                                                               (1,'Capítulo 1',22,currval('temporada_id_seq')),
                                                               (2,'Capítulo 2',22,currval('temporada_id_seq')),
                                                               (3,'Capítulo 3',22,currval('temporada_id_seq')),
                                                               (4,'Capítulo 4',22,currval('temporada_id_seq')),
                                                               (5,'Capítulo 5',22,currval('temporada_id_seq')),
                                                               (6,'Capítulo 6',22,currval('temporada_id_seq')),
                                                               (7,'Capítulo 7',22,currval('temporada_id_seq')),
                                                               (8,'Capítulo 8',22,currval('temporada_id_seq')),
                                                               (9,'Capítulo 9',22,currval('temporada_id_seq')),
                                                               (10,'Capítulo 10',22,currval('temporada_id_seq')),
                                                               (11,'Capítulo 11',22,currval('temporada_id_seq')),
                                                               (12,'Capítulo 12',22,currval('temporada_id_seq')),
                                                               (13,'Capítulo 13',22,currval('temporada_id_seq')),
                                                               (14,'Capítulo 14',22,currval('temporada_id_seq')),
                                                               (15,'Capítulo 15',22,currval('temporada_id_seq')),
                                                               (16,'Capítulo 16',22,currval('temporada_id_seq')),
                                                               (17,'Capítulo 17',22,currval('temporada_id_seq')),
                                                               (18,'Capítulo 18',22,currval('temporada_id_seq')),
                                                               (19,'Capítulo 19',22,currval('temporada_id_seq')),
                                                               (20,'Capítulo 20',22,currval('temporada_id_seq')),
                                                               (21,'Capítulo 21',22,currval('temporada_id_seq')),
                                                               (22,'Capítulo 22',22,currval('temporada_id_seq'));

INSERT INTO temporada (nro, id_serie) VALUES (3, 8);

INSERT INTO capitulo (nro, titulo, duracion, id_temporada) VALUES
                                                               (1,'Capítulo 1',22,currval('temporada_id_seq')),
                                                               (2,'Capítulo 2',22,currval('temporada_id_seq')),
                                                               (3,'Capítulo 3',22,currval('temporada_id_seq')),
                                                               (4,'Capítulo 4',22,currval('temporada_id_seq')),
                                                               (5,'Capítulo 5',22,currval('temporada_id_seq')),
                                                               (6,'Capítulo 6',22,currval('temporada_id_seq')),
                                                               (7,'Capítulo 7',22,currval('temporada_id_seq')),
                                                               (8,'Capítulo 8',22,currval('temporada_id_seq')),
                                                               (9,'Capítulo 9',22,currval('temporada_id_seq')),
                                                               (10,'Capítulo 10',22,currval('temporada_id_seq')),
                                                               (11,'Capítulo 11',22,currval('temporada_id_seq')),
                                                               (12,'Capítulo 12',22,currval('temporada_id_seq')),
                                                               (13,'Capítulo 13',22,currval('temporada_id_seq')),
                                                               (14,'Capítulo 14',22,currval('temporada_id_seq')),
                                                               (15,'Capítulo 15',22,currval('temporada_id_seq')),
                                                               (16,'Capítulo 16',22,currval('temporada_id_seq')),
                                                               (17,'Capítulo 17',22,currval('temporada_id_seq')),
                                                               (18,'Capítulo 18',22,currval('temporada_id_seq')),
                                                               (19,'Capítulo 19',22,currval('temporada_id_seq')),
                                                               (20,'Capítulo 20',22,currval('temporada_id_seq')),
                                                               (21,'Capítulo 21',22,currval('temporada_id_seq')),
                                                               (22,'Capítulo 22',22,currval('temporada_id_seq')),
                                                               (23,'Capítulo 23',22,currval('temporada_id_seq')),
                                                               (24,'Capítulo 24',22,currval('temporada_id_seq'));

/*Ginny & Georgia*/
INSERT INTO temporada (nro, id_serie) VALUES (1, 13);

INSERT INTO capitulo (nro, titulo, duracion, id_temporada) VALUES
                                                               (1,'Capítulo 1',55,currval('temporada_id_seq')),
                                                               (2,'Capítulo 2',55,currval('temporada_id_seq')),
                                                               (3,'Capítulo 3',55,currval('temporada_id_seq')),
                                                               (4,'Capítulo 4',55,currval('temporada_id_seq')),
                                                               (5,'Capítulo 5',55,currval('temporada_id_seq')),
                                                               (6,'Capítulo 6',55,currval('temporada_id_seq')),
                                                               (7,'Capítulo 7',55,currval('temporada_id_seq')),
                                                               (8,'Capítulo 8',55,currval('temporada_id_seq')),
                                                               (9,'Capítulo 9',55,currval('temporada_id_seq')),
                                                               (10,'Capítulo 10',55,currval('temporada_id_seq'));

INSERT INTO temporada (nro, id_serie) VALUES (2, 13);

INSERT INTO capitulo (nro, titulo, duracion, id_temporada) VALUES
                                                               (1,'Capítulo 1',55,currval('temporada_id_seq')),
                                                               (2,'Capítulo 2',55,currval('temporada_id_seq')),
                                                               (3,'Capítulo 3',55,currval('temporada_id_seq')),
                                                               (4,'Capítulo 4',55,currval('temporada_id_seq')),
                                                               (5,'Capítulo 5',55,currval('temporada_id_seq')),
                                                               (6,'Capítulo 6',55,currval('temporada_id_seq')),
                                                               (7,'Capítulo 7',55,currval('temporada_id_seq')),
                                                               (8,'Capítulo 8',55,currval('temporada_id_seq')),
                                                               (9,'Capítulo 9',55,currval('temporada_id_seq')),
                                                               (10,'Capítulo 10',55,currval('temporada_id_seq'));

INSERT INTO temporada (nro, id_serie) VALUES (3, 13);

INSERT INTO capitulo (nro, titulo, duracion, id_temporada) VALUES
                                                               (1,'Capítulo 1',55,currval('temporada_id_seq')),
                                                               (2,'Capítulo 2',55,currval('temporada_id_seq')),
                                                               (3,'Capítulo 3',55,currval('temporada_id_seq')),
                                                               (4,'Capítulo 4',55,currval('temporada_id_seq')),
                                                               (5,'Capítulo 5',55,currval('temporada_id_seq')),
                                                               (6,'Capítulo 6',55,currval('temporada_id_seq')),
                                                               (7,'Capítulo 7',55,currval('temporada_id_seq')),
                                                               (8,'Capítulo 8',55,currval('temporada_id_seq')),
                                                               (9,'Capítulo 9',55,currval('temporada_id_seq')),
                                                               (10,'Capítulo 10',55,currval('temporada_id_seq'));

/*Métodos de pago*/
INSERT INTO metodo_pago (descripcion) VALUES
                                          ('Tarjeta de crédito'),
                                          ('Transferencia bancaria'),
                                          ('PagoFácil'),
                                          ('Rapipago');

