#version 460
out vec4 frag_color;

const int MAX_RAY_STEPS = 128;
uniform vec2 uSize;
uniform float uTime;
uniform vec3 uPlayerPos;
const float PI = 3.14159265359;
uniform mat4 uInvView;
uniform mat4 uInvProj;

struct FlatNode {
    int child_indices[8]; // -1 means no child
    int is_leaf; // 1 for leaf, 0 otherwise
    // optional: int data; for voxel value
};

layout(std430, binding = 0) buffer OctreeBuffer {
    FlatNode nodes[];
};

bool getVoxel(ivec3 mapPos) {
    int size = WORLD_SIZE;
    ivec3 origin = ivec3(0);
    int index = 0; // root index

    while (true) {
        FlatNode node = nodes[index];

        // If we're at a leaf, voxel is present
        if (node.is_leaf == 1) return true;

        int half = size / 2;
        int child = 0;
        if (mapPos.x >= origin.x + half) child += 1;
        if (mapPos.y >= origin.y + half) child += 2;
        if (mapPos.z >= origin.z + half) child += 4;

        index = node.child_indices[child];
        if (index == -1) return false;

        // Update bounds
        if ((child & 1) != 0) origin.x += half;
        if ((child & 2) != 0) origin.y += half;
        if ((child & 4) != 0) origin.z += half;
        size = half;
    }
}

bvec3 dda(vec3 rayOrigin, vec3 rayDir) {
    ivec3 mapPos = ivec3(floor(rayOrigin));
    ivec3 rayStep = ivec3(sign(rayDir));
    vec3 deltaDist = abs(vec3(length(rayDir)) / rayDir);
    bvec3 mask;

    vec3 sideDist = (sign(rayDir) * (vec3(mapPos) - rayOrigin) + (sign(rayDir) * 0.5) + 0.5) * deltaDist;
    for (int i = 0; i < MAX_RAY_STEPS; i++) {
        if (getVoxel(mapPos)) break;
        mask = lessThanEqual(sideDist.xyz, min(sideDist.yzx, sideDist.zxy));
        sideDist += vec3(mask) * deltaDist;
        mapPos += ivec3(vec3(mask)) * rayStep;
        if (i == MAX_RAY_STEPS - 1) {
            mask = bvec3(0);
        }
    }
    return mask;
}

void main() {
    vec2 uv = (gl_FragCoord.xy * 2.0 - uSize) / uSize.y;
    vec3 rayOrigin = uPlayerPos;

    vec4 target = uInvProj * vec4(uv, 1, 1);
    vec3 rayDirection = (uInvView * vec4(normalize(target.xyz / target.w), 0)).xyz;
    vec3 lightDirection = -1 * normalize(vec3(sin(uTime * 2), 1, cos(uTime * 2)));

    // frag_color = vec4(raySphereIntersection(rayOrigin, rayDirection, .7, lightDirection), 1);
    // vec3 texCoord = vec3((uv.x + 1.0) / 2.0, (uv.y + 1.0) / 2.0, 0);
    // frag_color = vec4(vec3(texture(voxelMap, vec3(uv * 4, 0)) * 255), 1);
    bvec3 mask = dda(rayOrigin, rayDirection);
    if (mask.x) {
        frag_color.xyz = vec3(0.5);
    }
    else if (mask.y) {
        frag_color.xyz = vec3(1.0);
    }
    else if (mask.z) {
        frag_color.xyz = vec3(0.75);
    }
    else {
        frag_color.xyz = vec3(0.53, 0.81, 0.94);
    }
    frag_color.w = 1;
}
