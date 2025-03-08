## API

player.move_forward()
player.move_backward()
player.rotate_left()
player.rotate_right()
player.shoot()

asteroid.spawn()
asteroid.split() / asteroid.kill()

## Resources

- [How to handle delta time correctly](https://www.reddit.com/r/gamedev/comments/1embud0/how_to_handle_delta_time_correctly_in_multiplayer/)
- [Source Multiplayer Networking](https://developer.valvesoftware.com/wiki/Source_Multiplayer_Networking)

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
