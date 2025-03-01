import pygame
from circleshape import CircleShape
from constants import *

class Bullet(CircleShape):
    def __init__(self, x, y, velocity):
        super().__init__(x, y, BULLET_RADIUS)
        self.velocity = velocity

    def update(self, dt):
        self.position += BULLET_SPEED * dt * self.velocity

    def draw(self, screen):
        pygame.draw.circle(screen, "white", self.position, BULLET_RADIUS)
