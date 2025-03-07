import structlog
import pygame
import asteroids.game as game

log = structlog.get_logger()

def main():
    log.info("initializing pygame")
    pygame_num_success, pygame_num_failure = pygame.init()
    log.info("initialized pygame", num_success=pygame_num_success, num_failure=pygame_num_failure)

    with game.Game() as g:
        g.start()

if __name__ == "__main__":
    main()
