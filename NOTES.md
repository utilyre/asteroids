## API

player.move_forward()
player.move_backward()
player.rotate_left()
player.rotate_right()
player.shoot()

asteroid.spawn()
asteroid.split() / asteroid.kill()

---

game objects:

- player
- asteroid[]
- shot[]

states:

- player position
- asteroid[] position
- shot[] position

actions:

- player rotate right/left
- player move forward/backward
- asteroid spawn (server)

protocol:

message:
- version
- entity
- action
