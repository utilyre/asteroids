import structlog
import socket
import pygame
import asteroids.config as config
import asteroids.objects as objects

log = structlog.get_logger()

class Game:
    # connect to server
    # contain all state
    # init, start, exit

    def __init__(self):
        log.info("initializing game object")

        self.clock = pygame.time.Clock()
        self.dt = 0

        self.updatable = pygame.sprite.Group()
        self.drawable = pygame.sprite.Group()
        self.asteroids = pygame.sprite.Group()
        self.shots = pygame.sprite.Group()
        objects.Player.containers = (self.updatable, self.drawable) # TODO: How tf does this work?
        objects.Shot.containers = (self.updatable, self.drawable, self.shots)
        objects.Asteroid.containers = (self.updatable, self.drawable, self.asteroids)
        objects.AsteroidField.containers = (self.updatable)

        self.screen = None
        self.sock = None

        log.info("game object initialized")

    def __enter__(self):
        self.screen = pygame.display.set_mode((config.SCREEN_WIDTH, config.SCREEN_HEIGHT))
        log.info("set up display", width=config.SCREEN_WIDTH, height=config.SCREEN_HEIGHT)

        log.info("connecting to server via tcp")
        self.sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        self.sock.connect((config.SERVER_HOST, config.SERVER_PORT))
        log.info("successfully connected to server")

        return self

    def __exit__(self, exec_type, exec_value, traceback):
        log.info("closing tcp connection to the server")
        self.sock.close()

    def start(self):
        player = objects.Player(config.SCREEN_WIDTH / 2, config.SCREEN_HEIGHT / 2)
        asteroid_field = objects.AsteroidField(self.sock)
        log.info("initial game entities spawned")

        log.info("starting game loop")
        while True:
            for event in pygame.event.get():
                if event.type == pygame.QUIT:
                    return

            self.updatable.update(self.dt)

            for asteroid in self.asteroids:
                if player.collides_with(asteroid):
                    log.info("player died, game is over")
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
