#version 410 core

in vec3 FragColor;
out vec4 FragColorOut;

uniform vec3 color;

void main() {
    // Use uniform color override if it's not black
    if (color.x != 0.0 || color.y != 0.0 || color.z != 0.0) {
        FragColorOut = vec4(color, 1.0);
    } else {
        FragColorOut = vec4(FragColor, 1.0);
    }
}