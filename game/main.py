import logging
import pygame
import sys
from asteroid import Asteroid
from asteroidfield import AsteroidField
from constants import *
from player import Player
from shot import Shot
import socket
from network import send_message

def main():
    logging.basicConfig(level=logging.DEBUG)

    logging.info("initializing pygame")
    pygame_num_success, pygame_num_failure = pygame.init()
    logging.info("initialized pygame", extra={"num_success": pygame_num_success, "num_failure": pygame_num_failure})
    screen = pygame.display.set_mode((SCREEN_WIDTH, SCREEN_HEIGHT))
    logging.info("set up display", extra={"width": SCREEN_WIDTH, "height": SCREEN_HEIGHT})

    # all actions must be in player, right?
    # oh shit, asteroid field also has actions
    # but no, the server should spawn the asteroids

    with Game(screen) as g:
        g.start()

SERVER_HOST = "localhost"
SERVER_PORT = 3000

class Game:
    # connect to server
    # contain all state
    # init, start, exit

    def __init__(self, screen):
        logging.info("initializing game object")

        self.screen = screen

        self.clock = pygame.time.Clock()
        self.dt = 0

        self.updatable = pygame.sprite.Group()
        self.drawable = pygame.sprite.Group()
        self.asteroids = pygame.sprite.Group()
        self.shots = pygame.sprite.Group()
        Player.containers = (self.updatable, self.drawable) # TODO: How tf does this work?
        Shot.containers = (self.updatable, self.drawable, self.shots)
        Asteroid.containers = (self.updatable, self.drawable, self.asteroids)
        AsteroidField.containers = (self.updatable)

        self.sock = None

        logging.info("game object initialized")

    def __enter__(self):
        logging.info("connecting to server via tcp")
        self.sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        self.sock.connect((SERVER_HOST, SERVER_PORT))
        logging.info("successfully connected to server")
        return self

    def __exit__(self, exec_type, exec_value, traceback):
        logging.info("closing tcp connection to the server")
        self.sock.close()

    def start(self):
        player = Player(SCREEN_WIDTH / 2, SCREEN_HEIGHT / 2)
        asteroid_field = AsteroidField(self.sock)
        logging.info("initial game entities spawned")

        logging.info("starting game loop")
        while True:
            for event in pygame.event.get():
                if event.type == pygame.QUIT:
                    return

            self.updatable.update(self.dt)

            for asteroid in self.asteroids:
                if player.collides_with(asteroid):
                    logging.info("player died, game is over")
                    sys.exit()

            for shot in self.shots:
                for asteroid in self.asteroids:
                    if shot.collides_with(asteroid):
                        shot.kill()
                        asteroid.split()

            self.screen.fill("black")
            for sprite in self.drawable:
                sprite.draw(self.screen)
            pygame.display.flip()

            self.dt = self.clock.tick(60) / 1000

if __name__ == "__main__":
    main()
