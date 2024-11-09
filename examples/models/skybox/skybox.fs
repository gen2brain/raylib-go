#version 330

// input vertex attributes (from vertex shader)
in vec3 fragPosition;

// input uniform values
uniform samplerCube environmentMap;

// output fragment color
out vec4 finalColor;

void main()
{
    // fetch color from texture map
    vec3 color = texture(environmentMap, fragPosition).rgb;

    // calculate final fragment color
    finalColor = vec4(color, 1.0);
}