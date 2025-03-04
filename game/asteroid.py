import pygame
import random
from circleshape import CircleShape
from constants import *

class Asteroid(CircleShape):
    def __init__(self, x, y, radius):
        super().__init__(x, y, radius)

    def update(self, dt):
        self.position += dt * self.velocity

    def draw(self, screen):
        pygame.draw.circle(screen, "white", self.position, self.radius)

    def split(self):
        self.kill()
        if self.radius <= ASTEROID_MIN_RADIUS:
            return

        new_radius = self.radius - ASTEROID_MIN_RADIUS
        random_angle = random.uniform(20, 50)
        a1 = Asteroid(self.position.x, self.position.y, self.radius / 2)
        a2 = Asteroid(self.position.x, self.position.y, self.radius / 2)
        a1.velocity = 1.2 * self.velocity.rotate(-random_angle)
        a2.velocity = 1.2 * self.velocity.rotate(random_angle)
