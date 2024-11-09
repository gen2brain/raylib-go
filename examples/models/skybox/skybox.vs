#version 330

// input vertex attributes
in vec3 vertexPosition;

// input uniform values
uniform mat4 matProjection;
uniform mat4 matView;

// output vertex attributes (to fragment shader)
out vec3 fragPosition;

void main()
{
    // calculate fragment position based on model transformations
    fragPosition = vertexPosition;

    // remove translation from the view matrix
    mat4 rotView = mat4(mat3(matView));
    vec4 clipPos = matProjection*rotView*vec4(vertexPosition, 1.0);

    // calculate final vertex position
    gl_Position = clipPos;
}