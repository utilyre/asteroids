import pygame
import random
import asteroids.config as config
import asteroids.network as network

# Base class for game objects
class CircleShape(pygame.sprite.Sprite):
    def __init__(self, x, y, radius):
        # we will be using this later
        if hasattr(self, "containers"):
            super().__init__(self.containers)
        else:
            super().__init__()

        self.position = pygame.Vector2(x, y)
        self.velocity = pygame.Vector2(0, 0)
        self.radius = radius

    def draw(self, screen):
        # sub-classes must override
        pass

    def update(self, dt):
        # sub-classes must override
        pass

    def collides_with(self, other):
        l = self.position.distance_to(other.position)
        return l < self.radius + other.radius


class Shot(CircleShape):
    def __init__(self, x, y, velocity):
        super().__init__(x, y, config.SHOT_RADIUS)
        self.velocity = velocity

    def update(self, dt):
        self.position += dt * self.velocity

    def draw(self, screen):
        pygame.draw.circle(screen, "white", self.position, config.SHOT_RADIUS)

class Asteroid(CircleShape):
    def __init__(self, x, y, radius):
        super().__init__(x, y, radius)

    def update(self, dt):
        self.position += dt * self.velocity

    def draw(self, screen):
        pygame.draw.circle(screen, "white", self.position, self.radius)

    def split(self):
        self.kill()
        if self.radius <= config.ASTEROID_MIN_RADIUS:
            return

        new_radius = self.radius - config.ASTEROID_MIN_RADIUS
        random_angle = random.uniform(20, 50)
        a1 = Asteroid(self.position.x, self.position.y, self.radius / 2)
        a2 = Asteroid(self.position.x, self.position.y, self.radius / 2)
        a1.velocity = 1.2 * self.velocity.rotate(-random_angle)
        a2.velocity = 1.2 * self.velocity.rotate(random_angle)

class AsteroidField(pygame.sprite.Sprite):
    edges = [
        [
            pygame.Vector2(1, 0),
            lambda y: pygame.Vector2(-config.ASTEROID_MAX_RADIUS, y * config.SCREEN_HEIGHT),
        ],
        [
            pygame.Vector2(-1, 0),
            lambda y: pygame.Vector2(
                config.SCREEN_WIDTH + config.ASTEROID_MAX_RADIUS, y * config.SCREEN_HEIGHT
            ),
        ],
        [
            pygame.Vector2(0, 1),
            lambda x: pygame.Vector2(x * config.SCREEN_WIDTH, -config.ASTEROID_MAX_RADIUS),
        ],
        [
            pygame.Vector2(0, -1),
            lambda x: pygame.Vector2(
                x * config.SCREEN_WIDTH, config.SCREEN_HEIGHT + config.ASTEROID_MAX_RADIUS
            ),
        ],
    ]

    def __init__(self, sock):
        pygame.sprite.Sprite.__init__(self, self.containers)
        self.spawn_timer = 0.0
        self.sock = sock

    def spawn(self, radius, position, velocity):
        asteroid = Asteroid(position.x, position.y, radius)
        asteroid.velocity = velocity
        network.send_message(self.sock, 1, "asteroid/spawn", b"TODO")

    def update(self, dt):
        self.spawn_timer += dt
        if self.spawn_timer > config.ASTEROID_SPAWN_RATE:
            self.spawn_timer = 0

            # spawn a new asteroid at a random edge
            edge = random.choice(self.edges)
            speed = random.randint(40, 100)
            velocity = edge[0] * speed
            velocity = velocity.rotate(random.randint(-30, 30))
            position = edge[1](random.uniform(0, 1))
            kind = random.randint(1, config.ASTEROID_KINDS)
            self.spawn(config.ASTEROID_MIN_RADIUS * kind, position, velocity)

class Player(CircleShape):
    def __init__(self, x, y):
        super().__init__(x, y, config.PLAYER_RADIUS)
        self.rotation = 0
        self.cooldown = 0

    def draw(self, screen):
        pygame.draw.polygon(screen, "white", self.triangle(), 2)

    def update(self, dt):
        keys = pygame.key.get_pressed()

        # movement
        if keys[pygame.K_a]:
            self.rotate(-dt)
        if keys[pygame.K_d]:
            self.rotate(dt)
        if keys[pygame.K_w]:
            self.move(dt)
        if keys[pygame.K_s]:
            self.move(-dt)

        # shooting
        self.cooldown -= dt
        if keys[pygame.K_SPACE]:
            self.shoot()

    def shoot(self):
        if self.cooldown > 0:
            return

        Shot(
            self.position.x,
            self.position.y,
            config.PLAYER_SHOOT_SPEED * self.get_forward(),
        )
        self.cooldown = config.PLAYER_SHOOT_COOLDOWN

    def rotate(self, dt):
        self.rotation += config.PLAYER_TURN_SPEED * dt

    def move(self, dt):
        forward = self.get_forward()
        self.position += config.PLAYER_SPEED * dt * forward

    def get_forward(self):
        return pygame.Vector2(0, 1).rotate(self.rotation)

    def triangle(self):
        forward = pygame.Vector2(0, 1).rotate(self.rotation)
        right = pygame.Vector2(0, 1).rotate(self.rotation + 90) * self.radius / 1.5
        a = self.position + forward * self.radius
        b = self.position - forward * self.radius - right
        c = self.position - forward * self.radius + right
        return [a, b, c]

    def __repr__(self):
        return f"Player({self.position.x}, {self.position.y})"
