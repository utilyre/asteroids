import pygame
from constants import *

def main():
    print("Initializing pygame")
    success, failure = pygame.init()
    print(f"Successfully started {success} modules")
    if failure > 0:
        print(f"Failed to start {failure} modules")

    print("Setting up display")
    screen = pygame.display.set_mode((SCREEN_WIDTH, SCREEN_HEIGHT))
    print(f"Screen width: {SCREEN_WIDTH}")
    print(f"Screen height: {SCREEN_HEIGHT}")

    print("Starting Asteroids!")
    while True:
        for event in pygame.event.get():
            if event.type == pygame.QUIT:
                return

        screen.fill("black")
        pygame.display.flip()

if __name__ == "__main__":
    main()
