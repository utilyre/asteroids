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
    print("Creating a connection to the server")
    with socket.socket(socket.AF_INET, socket.SOCK_STREAM) as sock:
        print("Establishing the connection")
        sock.connect(("localhost", 3000))
        print("Connected to remote server")

        msg = {}
        msg["version"] = 1
        msg["entity"] = "player"
        msg["action"] = "move_forward"

        print("Sending test message to server")
        send_message(sock, msg)
        print("Test message sent to server")

    print("Initializing pygame")
    success, failure = pygame.init()
    print(f"Successfully started {success} modules")
    if failure > 0:
        print(f"Failed to start {failure} modules")

    print("Setting up display")
    screen = pygame.display.set_mode((SCREEN_WIDTH, SCREEN_HEIGHT))
    print(f"Screen width: {SCREEN_WIDTH}")
    print(f"Screen height: {SCREEN_HEIGHT}")

    updatable = pygame.sprite.Group()
    drawable = pygame.sprite.Group()
    asteroids = pygame.sprite.Group()
    shots = pygame.sprite.Group()
    Player.containers = (updatable, drawable) # TODO: How tf does this work?
    Shot.containers = (updatable, drawable, shots)
    Asteroid.containers = (updatable, drawable, asteroids)
    AsteroidField.containers = (updatable)

    player = Player(SCREEN_WIDTH / 2, SCREEN_HEIGHT / 2)
    asteroid_field = AsteroidField()

    # all actions must be in player, right?
    # oh shit, asteroid field also has actions
    # but no, the server should spawn the asteroids

    clock = pygame.time.Clock()
    dt = 0

    print("Starting Asteroids!")
    while True:
        for event in pygame.event.get():
            if event.type == pygame.QUIT:
                return

        updatable.update(dt)

        for asteroid in asteroids:
            if player.collides_with(asteroid):
                print("Game Over!")
                sys.exit()

        for shot in shots:
            for asteroid in asteroids:
                if shot.collides_with(asteroid):
                    shot.kill()
                    asteroid.split()

        screen.fill("black")
        for sprite in drawable:
            sprite.draw(screen)
        pygame.display.flip()

        dt = clock.tick(60) / 1000

if __name__ == "__main__":
    main()
