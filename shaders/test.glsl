//The raycasting code is somewhat based around a 2D raycasting toutorial found here:
//http://lodev.org/cgtutor/raycasting.html
#version 460

out vec4 fragColor;
const bool USE_BRANCHLESS_DDA = true;
const int MAX_RAY_STEPS = 64;
uniform float uTime;
uniform vec2 uSize;

float sdSphere(vec3 p, float d) {
    return length(p) - d;
}

float sdBox(vec3 p, vec3 b)
{
    vec3 d = abs(p) - b;
    return min(max(d.x, max(d.y, d.z)), 0.0) +
        length(max(d, 0.0));
}

bool getVoxel(ivec3 c) {
    vec3 p = vec3(c) + vec3(0.5);
    float d = min(max(-sdSphere(p, 7.5), sdBox(p, vec3(6.0))), -sdSphere(p, 25.0));
    return d < 0.0;
}

vec2 rotate2d(vec2 v, float a) {
    float sinA = sin(a);
    float cosA = cos(a);
    return vec2(v.x * cosA - v.y * sinA, v.y * cosA + v.x * sinA);
}

void main()
{
    vec2 screenPos = (gl_FragCoord.xy / uSize.xy) * 2.0 - 1.0;
    vec3 cameraDir = vec3(0.0, 0.0, 0.8);
    vec3 cameraPlaneU = vec3(1.0, 0.0, 0.0);
    vec3 cameraPlaneV = vec3(0.0, 1.0, 0.0) * uSize.y / uSize.x;
    vec3 rayDir = cameraDir + screenPos.x * cameraPlaneU + screenPos.y * cameraPlaneV;
    vec3 rayPos = vec3(0.0, 2.0 * sin(uTime * 2.7), -12.0);

    rayPos.xz = rotate2d(rayPos.xz, uTime);
    rayDir.xz = rotate2d(rayDir.xz, uTime);

    ivec3 mapPos = ivec3(floor(rayPos + 0.));

    vec3 deltaDist = abs(vec3(length(rayDir)) / rayDir);

    ivec3 rayStep = ivec3(sign(rayDir));

    vec3 sideDist = (sign(rayDir) * (vec3(mapPos) - rayPos) + (sign(rayDir) * 0.5) + 0.5) * deltaDist;

    bvec3 mask;

    for (int i = 0; i < MAX_RAY_STEPS; i++) {
        if (getVoxel(mapPos)) continue;
        if (USE_BRANCHLESS_DDA) {
            //Thanks kzy for the suggestion!
            mask = lessThanEqual(sideDist.xyz, min(sideDist.yzx, sideDist.zxy));
            /*bvec3 b1 = lessThan(sideDist.xyz, sideDist.yzx);
            	bvec3 b2 = lessThanEqual(sideDist.xyz, sideDist.zxy);
            	mask.x = b1.x && b2.x;
            	mask.y = b1.y && b2.y;
            	mask.z = b1.z && b2.z;*/
            //Would've done mask = b1 && b2 but the compiler is making me do it component wise.

            //All components of mask are false except for the corresponding largest component
            //of sideDist, which is the axis along which the ray should be incremented.

            sideDist += vec3(mask) * deltaDist;
            mapPos += ivec3(vec3(mask)) * rayStep;
        }
        else {
            if (sideDist.x < sideDist.y) {
                if (sideDist.x < sideDist.z) {
                    sideDist.x += deltaDist.x;
                    mapPos.x += rayStep.x;
                    mask = bvec3(true, false, false);
                }
                else {
                    sideDist.z += deltaDist.z;
                    mapPos.z += rayStep.z;
                    mask = bvec3(false, false, true);
                }
            }
            else {
                if (sideDist.y < sideDist.z) {
                    sideDist.y += deltaDist.y;
                    mapPos.y += rayStep.y;
                    mask = bvec3(false, true, false);
                }
                else {
                    sideDist.z += deltaDist.z;
                    mapPos.z += rayStep.z;
                    mask = bvec3(false, false, true);
                }
            }
        }
    }

    vec3 color;
    if (mask.x) {
        color = vec3(0.5);
    }
    if (mask.y) {
        color = vec3(1.0);
    }
    if (mask.z) {
        color = vec3(0.75);
    }
    fragColor.rgb = color;
    //fragColor.rgb = vec3(0.1 * noiseDeriv);
}
